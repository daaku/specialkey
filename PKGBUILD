# Builds specialkey from this checkout: makepkg here, then pacman -U the package.
pkgname=specialkey
pkgver=1
pkgrel=1
pkgdesc='Pick a special character from a popup grid, then type it and copy it'
arch=('x86_64')
url='https://github.com/daaku/specialkey'
license=('unknown') # the repository has no license yet
makedepends=('go' 'gtk4' 'gobject-introspection')
depends=(
  'cairo' 'gdk-pixbuf2' 'glib2' 'glibc' 'graphene' 'gtk4' 'harfbuzz' 'libgirepository' 'pango'
  'vulkan-icd-loader'
  'wl-clipboard' # wl-copy, so the clipboard outlives the popup
  'wtype'        # types the pick into the window that had focus
)

build() {
  cd ..
  export GOFLAGS='-buildmode=pie -trimpath -mod=readonly -modcacherw'

  # specialkey links GTK through cgo: makepkg's flags have to be handed to it.
  export CGO_CPPFLAGS="$CPPFLAGS" CGO_CFLAGS="$CFLAGS" CGO_CXXFLAGS="$CXXFLAGS" CGO_LDFLAGS="$LDFLAGS"
  go build -o build/specialkey .
}

package() {
  cd ..
  install -Dm755 build/specialkey "$pkgdir/usr/bin/specialkey"
  install -Dm644 deploy/specialkey.desktop "$pkgdir/usr/share/applications/specialkey.desktop"
}
