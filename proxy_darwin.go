package main

/*
#cgo LDFLAGS: -framework CoreFoundation -framework SystemConfiguration
#include <CoreFoundation/CoreFoundation.h>
#include <SystemConfiguration/SystemConfiguration.h>
#include <stdlib.h>

static CFDictionaryRef copyProxies(void) {
    return SCDynamicStoreCopyProxies(NULL);
}

static int dictInt(CFDictionaryRef dict, CFStringRef key) {
    CFNumberRef n = (CFNumberRef)CFDictionaryGetValue(dict, key);
    if (!n) return 0;
    int v = 0;
    CFNumberGetValue(n, kCFNumberIntType, &v);
    return v;
}

static char *cfStringToC(CFStringRef s) {
    if (!s) return NULL;
    CFIndex len = CFStringGetLength(s);
    CFIndex max = CFStringGetMaximumSizeForEncoding(len, kCFStringEncodingUTF8) + 1;
    char *buf = (char *)malloc(max);
    if (!buf) return NULL;
    if (!CFStringGetCString(s, buf, max, kCFStringEncodingUTF8)) {
        free(buf);
        return NULL;
    }
    return buf;
}

static char *dictString(CFDictionaryRef dict, CFStringRef key) {
    return cfStringToC((CFStringRef)CFDictionaryGetValue(dict, key));
}

static CFIndex arrayCount(CFDictionaryRef dict, CFStringRef key) {
    CFArrayRef a = (CFArrayRef)CFDictionaryGetValue(dict, key);
    if (!a) return 0;
    return CFArrayGetCount(a);
}

static char *arrayString(CFDictionaryRef dict, CFStringRef key, CFIndex i) {
    CFArrayRef a = (CFArrayRef)CFDictionaryGetValue(dict, key);
    if (!a) return NULL;
    return cfStringToC((CFStringRef)CFArrayGetValueAtIndex(a, i));
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func readConfig() (Config, error) {
	dict := C.copyProxies()
	if unsafe.Pointer(dict) == nil {
		return Config{}, nil
	}
	defer C.CFRelease(C.CFTypeRef(dict))

	proxy := func(enable, host, port C.CFStringRef) string {
		if C.dictInt(dict, enable) != 1 {
			return ""
		}
		cHost := C.dictString(dict, host)
		if cHost == nil {
			return ""
		}
		defer C.free(unsafe.Pointer(cHost))
		p := int(C.dictInt(dict, port))
		if p == 0 {
			return C.GoString(cHost)
		}
		return fmt.Sprintf("%s:%d", C.GoString(cHost), p)
	}

	cfg := Config{
		HTTP:  proxy(C.kSCPropNetProxiesHTTPEnable, C.kSCPropNetProxiesHTTPProxy, C.kSCPropNetProxiesHTTPPort),
		HTTPS: proxy(C.kSCPropNetProxiesHTTPSEnable, C.kSCPropNetProxiesHTTPSProxy, C.kSCPropNetProxiesHTTPSPort),
		FTP:   proxy(C.kSCPropNetProxiesFTPEnable, C.kSCPropNetProxiesFTPProxy, C.kSCPropNetProxiesFTPPort),
		SOCKS: proxy(C.kSCPropNetProxiesSOCKSEnable, C.kSCPropNetProxiesSOCKSProxy, C.kSCPropNetProxiesSOCKSPort),
	}

	n := C.arrayCount(dict, C.kSCPropNetProxiesExceptionsList)
	for i := C.CFIndex(0); i < n; i++ {
		cStr := C.arrayString(dict, C.kSCPropNetProxiesExceptionsList, i)
		if cStr == nil {
			continue
		}
		cfg.NoProxy = append(cfg.NoProxy, C.GoString(cStr))
		C.free(unsafe.Pointer(cStr))
	}

	return cfg, nil
}
