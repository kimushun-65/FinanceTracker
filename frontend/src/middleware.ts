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
 * Get authentication token from cookies
 */
function getAuthToken(request: NextRequest): string | null {
  return request.cookies.get('access_token')?.value || null;
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = getAuthToken(request);

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
    // Redirect to dashboard if already logged in and accessing home page
    if (pathname === '/' && token) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
    return NextResponse.next();
  }

  // Handle protected paths
  if (isProtectedPath(pathname)) {
    // Redirect to home page if no token
    if (!token) {
      const url = new URL('/', request.url);
      // Add query parameter to return to original URL after login
      url.searchParams.set('returnTo', pathname);
      return NextResponse.redirect(url);
    }

    // Continue processing if token exists
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