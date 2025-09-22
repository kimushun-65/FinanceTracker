import { NextRequest, NextResponse } from 'next/server';

// Authentication required paths
const PROTECTED_PATHS = [
  '/dashboard',
  '/transactions',
  '/reports',
  '/accounts',
  '/budget',
  '/settings',
];

// Public paths that don't require authentication
const PUBLIC_PATHS = [
  '/',
  '/login',
  '/callback',
  '/api/auth',
];

/**
 * Check if the given path is protected
 */
function isProtectedPath(pathname: string): boolean {
  return PROTECTED_PATHS.some(path => pathname.startsWith(path));
}

/**
 * Check if the given path is public
 */
function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATHS.some(path => pathname === path || pathname.startsWith(path));
}

/**
 * Check if access token exists in cookies (simplified check for middleware)
 */
function hasAccessToken(request: NextRequest): boolean {
  const token = request.cookies.get('access_token');
  return !!(token && token.value && token.value.length > 10);
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Skip processing for static files, Next.js internal paths, and API routes
  if (
    pathname.startsWith('/_next') ||
    pathname.startsWith('/api') ||
    pathname.includes('.')
  ) {
    return NextResponse.next();
  }

  // Skip authentication check for public paths
  if (isPublicPath(pathname)) {
    // Note: Removed automatic redirect from / to /dashboard to prevent loops
    // Let the frontend handle the login flow
    return NextResponse.next();
  }

  // Handle protected paths
  if (isProtectedPath(pathname)) {
    // Check if access token exists in cookies
    const hasToken = hasAccessToken(request);
    if (!hasToken) {
      const url = new URL('/', request.url);
      // Add query parameter to return to original URL after login
      url.searchParams.set('returnTo', pathname);
      return NextResponse.redirect(url);
    }

    // Continue processing if token exists
    // Note: Actual token validation will be done by individual page components
    return NextResponse.next();
  }

  // Continue processing for other paths
  return NextResponse.next();
}

// Configure paths where middleware should be applied
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - images, icons, etc.
     */
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.).*)',
  ],
};