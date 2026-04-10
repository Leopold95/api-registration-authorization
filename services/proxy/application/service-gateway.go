package application

import (
	"api-registration-authorization/services/proxy/domain"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3/log"
)

type ServiceEndpoint struct {
	Host      string
	Port      string
	LastCheck time.Time
	IsPrivate bool
}

type GatewayService struct {
	routes map[string]*ServiceEndpoint
	mu     sync.RWMutex
}

func InitGateway() *GatewayService {
	gateway := &GatewayService{
		routes: make(map[string]*ServiceEndpoint),
	}

	separator := "~"

	auth := strings.Split(os.Getenv(domain.ServiceAuthorization), separator)
	register := strings.Split(os.Getenv(domain.ServiceRegistration), separator)
	users := strings.Split(os.Getenv(domain.ServiceUsers), separator)
	shops := strings.Split(os.Getenv(domain.ServiceShops), separator)
	requsts := strings.Split(os.Getenv(domain.ServiceRequests), separator)
	commnents := strings.Split(os.Getenv(domain.ServiceComments), separator)

	log.Debug(auth)
	log.Debug(register)
	log.Debug(users)
	log.Debug(shops)
	log.Debug(requsts)
	log.Debug(commnents)

	requireAuth, _ := strconv.ParseBool(auth[3])
	gateway.Register(auth[0], auth[1], auth[2], requireAuth)

	requireRegister, _ := strconv.ParseBool(register[3])
	gateway.Register(register[0], register[1], register[2], requireRegister)

	requireUsers, _ := strconv.ParseBool(users[3])
	gateway.Register(users[0], users[1], users[2], requireUsers)

	requireShops, _ := strconv.ParseBool(shops[3])
	gateway.Register(shops[0], shops[1], shops[2], requireShops)

	requireRequsts, _ := strconv.ParseBool(requsts[3])
	gateway.Register(requsts[0], requsts[1], requsts[2], requireRequsts)

	requireComments, _ := strconv.ParseBool(commnents[3])
	gateway.Register(commnents[0], commnents[1], commnents[2], requireComments)

	//gateway.Register("/api/v2/auth/authorize", "http://127.0.0.1", "3000", true)
	//gateway.Register("/api/v2/auth/register", "http://127.0.0.1", "3001", true)
	//
	//gateway.Register("/api/v1/user", "http://host.docker.internal", "6012", false)
	//gateway.Register("/api/shops", "http://host.docker.internal", "6013", false)
	//gateway.Register("/api/requests", "http://host.docker.internal", "6014", false)
	//gateway.Register("/api/v1/comments", "http://host.docker.internal", "6015", false)

	return gateway
}

// Register a new service endpoint
func (g *GatewayService) Register(path, host, port string, isPrivate bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.routes[path] = &ServiceEndpoint{
		Host:      host,
		Port:      port,
		LastCheck: time.Time{},
		IsPrivate: isPrivate,
	}
}

func (g *GatewayService) GetEndpoint(path string) (*ServiceEndpoint, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	endpoint, exists := g.routes[path]
	if !exists {
		return nil, false
	}
	return endpoint, true
}

func (g *GatewayService) GetRoutes() map[string]*ServiceEndpoint {
	return g.routes
}
