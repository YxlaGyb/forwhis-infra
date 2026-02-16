package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIss        = "lmoadll"
	jwtAud        = "lmoadll"
	keySize       = 32
	rotationDays  = 7
	maxActiveKeys = 3
	tokenTTL      = 24 * time.Hour
	localCacheTTL = time.Minute
	retryMax      = 3
	retryWait     = 100 * time.Millisecond
)
const (
	keyHashKey       = "jwt:keys"
	keyMetaHashKey   = "jwt:keys:meta"
	activeKeyKey     = "jwt:active_key"
	blacklistPrefix  = "jwt:blacklist:"
	userTokensPrefix = "jwt:user:"
	lockKey          = "jwt:lock:rotate"
	cacheChannel     = "jwt:cache:invalidate"
)

var (
	ErrMissingKid       = errors.New("missing kid in header")
	ErrInvalidKid       = errors.New("invalid kid")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenRevoked     = errors.New("token revoked")
	ErrInvalidIssuer    = errors.New("invalid issuer")
	ErrInvalidAudience  = errors.New("invalid audience")
	ErrRedisUnavailable = errors.New("redis service unavailable")
)

type Key struct {
	ID        string
	Secret    []byte
	CreatedAt time.Time
}
type JWTManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex

	keys     map[string]*Key
	activeID string
	expired  time.Time
}
type CustomClaims struct {
	JID   string `json:"jid"`
	UID   int    `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var (
	manager *JWTManager
	once    sync.Once
)

func Init() {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		manager = &JWTManager{
			ctx:    ctx,
			cancel: cancel,
			keys:   make(map[string]*Key),
		}
		manager.wg.Add(2)
		go manager.rotateLoop()
		go manager.subscribeLoop()
	})
}

func Shutdown() {
	if manager != nil {
		manager.cancel()
		manager.wg.Wait()
	}
}

func ensureInit() {
	if manager == nil {
		Init()
	}
}

// --- Key 轮转与同步 ---

func (m *JWTManager) rotateLoop() {
	defer m.wg.Done()

	// 启动时检查
	time.Sleep(2 * time.Second)
	if m.getActiveKey() == nil {
		_ = m.rotateKey()
	}

	ticker := time.NewTicker(rotationDays * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			_ = m.rotateKey()
		}
	}
}

func (m *JWTManager) subscribeLoop() {
	defer m.wg.Done()
	pubsub := db.Client().Subscribe(m.ctx, cacheChannel)
	defer pubsub.Close()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-pubsub.Channel():
			time.Sleep(jitter()) // 防止缓存雪崩
			m.mu.Lock()
			m.keys = make(map[string]*Key) // 清空本地缓存
			m.activeID = ""
			m.mu.Unlock()
		}
	}
}

func (m *JWTManager) rotateKey() error {
	client := db.Client()

	// 1. 获取分布式锁
	acquired := false
	for i := 0; i < retryMax; i++ {
		if ok, _ := client.SetNX(m.ctx, lockKey, "1", 10*time.Second).Result(); ok {
			acquired = true
			break
		}
		time.Sleep(retryWait)
	}

	if !acquired {
		if key := m.getActiveKey(); key != nil {
			return nil
		}
		return ErrRedisUnavailable
	}
	defer client.Del(m.ctx, lockKey)

	// 2. 双重检查
	if id, _ := client.Get(m.ctx, activeKeyKey).Result(); id != "" {
		return nil
	}

	// 3. 生成并存储新 Key
	newKey := &Key{
		ID:        generateID(),
		Secret:    generateSecret(),
		CreatedAt: time.Now(),
	}

	pipe := client.Pipeline()
	pipe.HSet(m.ctx, keyHashKey, newKey.ID, string(newKey.Secret))
	pipe.HSet(m.ctx, keyMetaHashKey, newKey.ID, newKey.CreatedAt.Format(time.RFC3339))
	pipe.Set(m.ctx, activeKeyKey, newKey.ID, 0)
	if _, err := pipe.Exec(m.ctx); err != nil {
		return err
	}

	// 4. 异步清理与通知
	go m.cleanupOldKeys(context.Background(), newKey.ID)

	m.mu.Lock()
	m.keys[newKey.ID] = newKey
	m.activeID = newKey.ID
	m.expired = time.Now().Add(localCacheTTL + jitter())
	m.mu.Unlock()

	client.Publish(m.ctx, cacheChannel, "invalidate")
	return nil
}

func (m *JWTManager) cleanupOldKeys(ctx context.Context, newID string) {
	client := db.Client()
	n, _ := client.HLen(ctx, keyHashKey).Result()
	if n <= maxActiveKeys {
		return
	}

	meta, _ := client.HGetAll(ctx, keyMetaHashKey).Result()
	type keyMeta struct {
		id string
		t  time.Time
	}
	keys := make([]keyMeta, 0, len(meta))

	for kid, timeStr := range meta {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
			keys = append(keys, keyMeta{kid, t})
		}
	}

	// 按时间升序排序（旧的在前）
	sort.Slice(keys, func(i, j int) bool { return keys[i].t.Before(keys[j].t) })

	// 删除最旧的 Key，保留 maxActiveKeys 个
	for i := 0; i < len(keys)-maxActiveKeys; i++ {
		if keys[i].id != newID {
			client.HDel(ctx, keyHashKey, keys[i].id)
			client.HDel(ctx, keyMetaHashKey, keys[i].id)
		}
	}
}

func (m *JWTManager) getActiveKey() *Key {
	m.mu.RLock()
	if time.Now().Before(m.expired) && m.activeID != "" {
		if k := m.keys[m.activeID]; k != nil {
			m.mu.RUnlock()
			return k
		}
	}
	m.mu.RUnlock()

	// Cache Miss
	id, err := db.Client().Get(m.ctx, activeKeyKey).Result()
	if err != nil || id == "" {
		return nil
	}
	return m.getKey(id)
}

func (m *JWTManager) getKey(kid string) *Key {
	m.mu.RLock()
	if k, ok := m.keys[kid]; ok && time.Now().Before(m.expired) {
		m.mu.RUnlock()
		return k
	}
	m.mu.RUnlock()

	secret, err := db.Client().HGet(m.ctx, keyHashKey, kid).Result()
	if err != nil || secret == "" {
		return nil
	}

	k := &Key{ID: kid, Secret: []byte(secret)}
	m.mu.Lock()
	m.keys[kid] = k
	m.mu.Unlock()
	return k
}

// --- Token 接口 ---

func GenerateToken(jid string, uid int, email string) (string, error) {
	ensureInit()
	return manager.generateToken(jid, uid, email)
}

func (m *JWTManager) generateToken(jid string, uid int, email string) (string, error) {
	key := m.getActiveKey()
	if key == nil {
		if err := m.rotateKey(); err != nil {
			return "", err
		}
		if key = m.getActiveKey(); key == nil {
			return "", ErrRedisUnavailable
		}
	}

	now := time.Now()
	claims := CustomClaims{
		JID:   jid,
		UID:   uid,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIss,
			Audience:  jwt.ClaimStrings{jwtAud},
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = key.ID

	signed, err := token.SignedString(key.Secret)
	if err != nil {
		return "", err
	}

	// 异步记录：优化接口响应速度
	go func(uid int, jid string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		key := userTokensPrefix + strconv.Itoa(uid)
		db.SAdd(ctx, key, jid)
		db.Expire(ctx, key, tokenTTL)
	}(uid, jid)

	return signed, nil
}

func ValidateToken(tokenString string) (*types.Payload, error) {
	ensureInit()
	return manager.validateToken(tokenString)
}

func (m *JWTManager) validateToken(tokenString string) (*types.Payload, error) {
	var claims CustomClaims
	// ParseWithClaims 会自动验证过期时间 (exp)
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, ErrMissingKid
		}

		key := m.getKey(kid)
		if key == nil {
			return nil, ErrInvalidKid
		}
		return key.Secret, nil
	})

	if err != nil {
		return nil, err
	}
	if isRevoked(claims.JID) {
		return nil, ErrTokenRevoked
	}

	// 业务级校验
	if claims.Issuer != jwtIss {
		return nil, ErrInvalidIssuer
	}

	// 安全修复：防止Audience为空导致Panic
	audMatch := false
	for _, aud := range claims.Audience {
		if aud == jwtAud {
			audMatch = true
			break
		}
	}
	if !audMatch {
		return nil, ErrInvalidAudience
	}

	// 构造 Payload 返回给上层业务
	return &types.Payload{
		JID:     claims.JID,
		UID:     claims.UID,
		Email:   claims.Email,
		Expired: claims.ExpiresAt.Unix(),
		ISS:     claims.Issuer,
		AUD:     jwtAud,
	}, nil
}

// --- 辅助工具 ---

func ValidateTokenFromContext(c *gin.Context) (*types.Payload, error) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	} else {
		return nil, ErrInvalidToken
	}
	return ValidateToken(token)
}

func RevokeUserTokens(uid int) {
	ensureInit()
	uidStr := strconv.Itoa(uid)
	ctx := context.Background()
	client := db.Client()

	tokens, err := client.SMembers(ctx, userTokensPrefix+uidStr).Result()
	if err != nil || len(tokens) == 0 {
		return
	}

	pipe := client.Pipeline()
	for _, jid := range tokens {
		pipe.Set(ctx, blacklistPrefix+jid, "1", tokenTTL)
	}
	pipe.Del(ctx, userTokensPrefix+uidStr)
	pipe.Exec(ctx)
}

func isRevoked(jid string) bool {
	n, _ := db.Exists(blacklistPrefix + jid)
	return n > 0
}

func GetPayload(c *gin.Context) *types.Payload {
	if p, ok := c.Get("payload"); ok {
		return p.(*types.Payload)
	}
	return nil
}

// 随机抖动，避免缓存同时失效
func jitter() time.Duration {
	n, _ := rand.Int(rand.Reader, big.NewInt(5000)) // 0-5s
	return time.Duration(n.Int64()) * time.Millisecond
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSecret() []byte {
	b := make([]byte, keySize)
	rand.Read(b)
	return b
}
