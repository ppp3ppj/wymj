package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
    envMap, err := godotenv.Read(path)
    if err != nil {
        log.Fatalf("load env file error: %v", err)
    }
    return &config{
        app: &app{
            host: envMap["APP_HOST"],
            port: func() int {
                p, err := strconv.Atoi(envMap["APP_PORT"])
                if err != nil {
                    log.Fatalf("convert port to int error: %v", err)
                }
                return p
            }(),
            name: envMap["APP_NAME"],
            version: envMap["APP_VERSION"],
            readTimeout: func() time.Duration {
                t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
                if err != nil {
                    log.Fatalf("convert readTimeout to int error: %v", err)
                }
                return time.Duration(int64(t) * int64(math.Pow10(9)))
            }(),
            writeTimeout: func() time.Duration {
                t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
                if err != nil {
                    log.Fatalf("convert writeTimeout to int error: %v", err)
                }
                return time.Duration(int64(t) * int64(math.Pow10(9)))
            }(),
            bodyLimit: func() int {
                b, err := strconv.Atoi(envMap["APP_PORT"])
                if err != nil {
                    log.Fatalf("convert bodyLimit to int error: %v", err)
                }
                return b
            }(),
            fileLimit: func() int {
                f, err := strconv.Atoi(envMap["APP_PORT"])
                if err != nil {
                    log.Fatalf("convert fileLimit to int error: %v", err)
                }
                return f
            }(),
            gcpbucket:  envMap["APP_GCP_BUCKET"],
        },
        db: &db{
            host: envMap["DB_HOST"],
            port: func() int {
                p, err := strconv.Atoi(envMap["DB_PORT"])
                if err != nil {
                    log.Fatalf("convert port to int error: %v", err)
                }
                return p
            }(),
            protocol: envMap["DB_PROTOCOL"],
            username: envMap["DB_USERNAME"],
            password: envMap["DB_PASSWORD"],
            database: envMap["DB_DATABASE"],
            sslmode: envMap["DB_SSL_MODE"],
            maxConnections: func() int {
                m, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
                if err != nil {
                    log.Fatalf("convert maxConnections to int error: %v", err)
                }
                return m
            }(),
        },
        jwt: &jwt{
            adminKey: envMap["JWT_SECRET_KEY"],
            secretKey: envMap["JWT_ADMIN_KEY"],
            apiKey: envMap["JWT_API_KEY"],
            accessExpireAt: func() int {
                a, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
                if err != nil {
                    log.Fatalf("convert accessExpireAt to int error: %v", err)
                }
                return a
            }(),
            refreshExpireAt: func() int {
                r, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
                if err != nil {
                    log.Fatalf("convert refreshExpireAt to int error: %v", err)
                }
                return r
            }(),
        },
    }
}

type IConfig interface {
    App() IAppconfig
    Db() IDbconfig
    Jwt() IJwtconfig
}

type config struct {
    app *app
    db *db
    jwt *jwt
}

type IAppconfig interface {
    // host:port
    Url() string
    Name() string
    Version() string
    ReadTimeout() time.Duration
    WriteTimeout() time.Duration
    BodyLimit() int
    FileLimit() int
    Gcpbucket() string
}

type app struct {
    host string
    port int
    name string
    version string
    readTimeout time.Duration
    writeTimeout time.Duration
    bodyLimit int //in bytes
    fileLimit int //in bytes
    gcpbucket string
}

func (c *config) App() IAppconfig {
    return c.app
}

func (a *app) Url() string { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string { return a.name }
func (a *app) Version() string { return a.version }
func (a *app) ReadTimeout() time.Duration { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int { return a.bodyLimit }
func (a *app) FileLimit() int { return a.fileLimit }
func (a *app) Gcpbucket() string { return a.gcpbucket }

type IDbconfig interface {
    Url() string
    MaxConnections() int
}

type db struct {
    host string
    port int
    protocol string
    username string
    password string
    database string
    sslmode string
    maxConnections int
}

func (c *config) Db() IDbconfig {
    return c.db
}

func (d *db) Url() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        d.host,
        d.port,
        d.username,
        d.password,
        d.database,
        d.sslmode,
    )
}
func (d *db) MaxConnections() int { return d.maxConnections }

type IJwtconfig interface {
    SecretKey() []byte
    AdminKey() []byte
    ApiKey() []byte
    AccessExpireAt() int
    RefreshExpireAt() int

    SetJwtAccessExpireAt(t int)
    SetJwtExpireAt(t int)
}

type jwt struct {
    adminKey string
    secretKey string
    apiKey string
    accessExpireAt int //in seconds
    refreshExpireAt int //in seconds
}

func (c *config) Jwt() IJwtconfig {
    return c.jwt
}

func (j *jwt) SecretKey() []byte { return []byte(j.secretKey) }
func (j *jwt) AdminKey() []byte { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte { return []byte(j.apiKey) }
func (j *jwt) AccessExpireAt() int { return j.accessExpireAt }
func (j *jwt) RefreshExpireAt() int { return j.refreshExpireAt }
func (j *jwt) SetJwtAccessExpireAt(t int) { j.accessExpireAt = t }
func (j *jwt) SetJwtExpireAt(t int) { j.refreshExpireAt = t }
