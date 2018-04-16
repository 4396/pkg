package dep

import (
	"go/build"
	"os"
	"strings"
)

var stdPackages = map[string]bool{
	"C":                         true,
	"builtin":                   true,
	"archive/tar":               true,
	"archive/zip":               true,
	"bufio":                     true,
	"bytes":                     true,
	"compress/bzip2":            true,
	"compress/flate":            true,
	"compress/gzip":             true,
	"compress/lzw":              true,
	"compress/zlib":             true,
	"container/heap":            true,
	"container/list":            true,
	"container/ring":            true,
	"context":                   true,
	"crypto":                    true,
	"crypto/aes":                true,
	"crypto/cipher":             true,
	"crypto/des":                true,
	"crypto/dsa":                true,
	"crypto/ecdsa":              true,
	"crypto/elliptic":           true,
	"crypto/hmac":               true,
	"crypto/md5":                true,
	"crypto/rand":               true,
	"crypto/rc4":                true,
	"crypto/rsa":                true,
	"crypto/sha1":               true,
	"crypto/sha256":             true,
	"crypto/sha512":             true,
	"crypto/subtle":             true,
	"crypto/tls":                true,
	"crypto/x509":               true,
	"crypto/x509/pkix":          true,
	"database/sql":              true,
	"database/sql/driver":       true,
	"debug/dwarf":               true,
	"debug/elf":                 true,
	"debug/gosym":               true,
	"debug/macho":               true,
	"debug/pe":                  true,
	"debug/plan9obj":            true,
	"encoding":                  true,
	"encoding/ascii85":          true,
	"encoding/asn1":             true,
	"encoding/base32":           true,
	"encoding/base64":           true,
	"encoding/binary":           true,
	"encoding/csv":              true,
	"encoding/gob":              true,
	"encoding/hex":              true,
	"encoding/json":             true,
	"encoding/pem":              true,
	"encoding/xml":              true,
	"errors":                    true,
	"expvar":                    true,
	"flag":                      true,
	"fmt":                       true,
	"go/ast":                    true,
	"go/build":                  true,
	"go/constant":               true,
	"go/doc":                    true,
	"go/format":                 true,
	"go/importer":               true,
	"go/internal/gccgoimporter": true,
	"go/internal/gcimporter":    true,
	"go/parser":                 true,
	"go/printer":                true,
	"go/scanner":                true,
	"go/token":                  true,
	"go/types":                  true,
	"hash":                      true,
	"hash/adler32":              true,
	"hash/crc32":                true,
	"hash/crc64":                true,
	"hash/fnv":                  true,
	"html":                      true,
	"html/template":             true,
	"image":                     true,
	"image/color":               true,
	"image/color/palette":       true,
	"image/draw":                true,
	"image/gif":                 true,
	"image/internal/imageutil":  true,
	"image/jpeg":                true,
	"image/png":                 true,
	"index/suffixarray":         true,
	"internal/race":             true,
	"internal/singleflight":     true,
	"internal/testenv":          true,
	"internal/trace":            true,
	"io":                        true,
	"io/ioutil":                 true,
	"log":                       true,
	"log/syslog":                true,
	"math":                      true,
	"math/big":                  true,
	"math/cmplx":                true,
	"math/rand":                 true,
	"mime":                      true,
	"mime/multipart":            true,
	"mime/quotedprintable":      true,
	"net":                     true,
	"net/http":                true,
	"net/http/cgi":            true,
	"net/http/cookiejar":      true,
	"net/http/fcgi":           true,
	"net/http/httptest":       true,
	"net/http/httptrace":      true,
	"net/http/httputil":       true,
	"net/http/internal":       true,
	"net/http/pprof":          true,
	"net/internal/socktest":   true,
	"net/mail":                true,
	"net/rpc":                 true,
	"net/rpc/jsonrpc":         true,
	"net/smtp":                true,
	"net/textproto":           true,
	"net/url":                 true,
	"os":                      true,
	"os/exec":                 true,
	"os/signal":               true,
	"os/user":                 true,
	"path":                    true,
	"path/filepath":           true,
	"reflect":                 true,
	"regexp":                  true,
	"regexp/syntax":           true,
	"runtime":                 true,
	"runtime/cgo":             true,
	"runtime/debug":           true,
	"runtime/internal/atomic": true,
	"runtime/internal/sys":    true,
	"runtime/pprof":           true,
	"runtime/race":            true,
	"runtime/trace":           true,
	"sort":                    true,
	"strconv":                 true,
	"strings":                 true,
	"sync":                    true,
	"sync/atomic":             true,
	"syscall":                 true,
	"testing":                 true,
	"testing/iotest":          true,
	"testing/quick":           true,
	"text/scanner":            true,
	"text/tabwriter":          true,
	"text/template":           true,
	"text/template/parse":     true,
	"time":                    true,
	"unicode":                 true,
	"unicode/utf16":           true,
	"unicode/utf8":            true,
	"unsafe":                  true,
}

func listImports(ctx build.Context, dones map[string]bool, path, root, src string, test bool) (imports []string, err error) {
	_, ok := dones[path]
	if ok {
		return
	}
	dones[path] = true

	pkg, err := ctx.Import(path, src, build.AllowBinary)
	if err != nil {
		if _, ok := err.(*build.NoGoError); ok {
			err = nil
		}
		return
	}

	paths := pkg.Imports
	if test {
		paths = append(paths, pkg.TestImports...)
	}

	for _, path := range paths {
		_, ok := stdPackages[path]
		if ok {
			continue
		}

		if strings.HasPrefix(path, root) {
			var subImports []string
			subImports, err = listImports(ctx, dones, path, root, src, test)
			if err != nil {
				return
			}
			imports = append(imports, subImports...)
		} else {
			imports = append(imports, path)
		}
	}
	return
}

func Imports(path, src string, test bool) (imports []string, err error) {
	ctx := build.Default
	gopath := os.Getenv("GOPATH")
	ctx.GOPATH = src + string(os.PathListSeparator) + gopath
	dones := make(map[string]bool)
	imports, err = listImports(ctx, dones, path, path, src, test)
	return
}
