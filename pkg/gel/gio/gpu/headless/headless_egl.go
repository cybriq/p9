// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || windows || openbsd
// +build linux freebsd windows openbsd

package headless

import (
	"github.com/cybriq/p9/pkg/gel/gio/internal/egl"
)

func newGLContext() (context, error) {
	return egl.NewContext(egl.EGL_DEFAULT_DISPLAY)
}
