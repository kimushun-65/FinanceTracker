/**
 * Cookie操作のユーティリティ関数
 */

/**
 * Cookieを取得する
 * @param name Cookie名
 * @returns Cookie値（存在しない場合はnull）
 */
export function getCookie(name: string): string | null {
    if (typeof window === 'undefined') return null;
  
    const nameEQ = `${name}=`;
    const cookies = document.cookie.split(';');
  
    for (let cookie of cookies) {
      cookie = cookie.trim();
      if (cookie.indexOf(nameEQ) === 0) {
        return decodeURIComponent(cookie.substring(nameEQ.length));
      }
    }
    return null;
  }
  
  /**
   * 認証トークン用の定数
   */
  export const AUTH_TOKEN_KEY = 'access_token';