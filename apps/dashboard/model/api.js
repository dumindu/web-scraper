// Base API client with authentication handling
const API_BASE_URL = 'http://localhost:8080/v1'

class ApiClient {
  constructor() {
    this.baseUrl = API_BASE_URL
  }

  getToken() {
    return localStorage.getItem('access_token')
  }

  setToken(token) {
    localStorage.setItem('access_token', token)
  }

  clearToken() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`
    
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers
    }

    // Always get fresh token from localStorage
    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    } else {
      // No token available
    }

    const config = {
      ...options,
      headers
    }

    try {
      const response = await fetch(url, config)
      
      if (response.status === 401) {
        // Check if this is a login request - don't redirect for login failures
        if (endpoint === '/users/sign-in') {
          // Let login errors fall through to normal error handling
        } else {
          // Token expired, clear auth and redirect
          this.clearToken()
          if (window.navigateTo) {
            window.navigateTo('/login')
          } else {
            window.location.href = '/login'
          }
          return
        }
      }

      if (!response.ok) {
        let errorMessage = 'API request failed'
        try {
          const contentType = response.headers.get('content-type')
          if (contentType && contentType.includes('application/json')) {
            const error = await response.json()
            
            // Handle 400 validation errors with errors array
            if (response.status === 400 && error.errors && Array.isArray(error.errors)) {
              errorMessage = error.errors.join(', ')
            } else {
              errorMessage = error.error || error.message || errorMessage
            }
          } else {
            // If not JSON, use response text or status
            const text = await response.text()
            errorMessage = text || response.statusText || errorMessage
          }
        } catch (e) {
          // If parsing fails, use the response status text
          errorMessage = response.statusText || errorMessage
        }
        // Include status code in error message
        errorMessage = `${response.status}: ${errorMessage}`
        throw new Error(errorMessage)
      }

      // Handle successful responses
      const contentType = response.headers.get('content-type')
      
      // Read response as text first
      const textResponse = await response.text()
      
      // Handle empty responses
      if (!textResponse || textResponse.trim() === '') {
        return { success: true, message: 'Request completed successfully' }
      }
      
      if (contentType && contentType.includes('application/json')) {
        try {
          const jsonResponse = JSON.parse(textResponse)
          return jsonResponse
        } catch (jsonError) {
          // If JSON parsing fails but response was successful, return success
          return { success: true, message: 'Request completed successfully' }
        }
      } else {
        // Return text for non-JSON responses
        return textResponse
      }
    } catch (error) {
      throw error
    }
  }

  async get(endpoint) {
    return this.request(endpoint, { method: 'GET' })
  }

  async post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  async postFormData(endpoint, formData) {
    const url = `${this.baseUrl}${endpoint}`
    
    const headers = {}
    
    // Always get fresh token from localStorage
    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    // Don't set Content-Type for FormData, let browser set it with boundary
    const config = {
      method: 'POST',
      headers,
      body: formData
    }

    try {
      const response = await fetch(url, config)
      
      if (response.status === 401) {
        // Check if this is a login request - don't redirect for login failures
        if (endpoint === '/users/sign-in') {
          // Let login errors fall through to normal error handling
        } else {
          // Token expired, clear auth and redirect
          this.clearToken()
          if (window.navigateTo) {
            window.navigateTo('/login')
          } else {
            window.location.href = '/login'
          }
          return
        }
      }

      if (!response.ok) {
        let errorMessage = 'API request failed'
        try {
          const contentType = response.headers.get('content-type')
          if (contentType && contentType.includes('application/json')) {
            const error = await response.json()
            
            // Handle 400 validation errors with errors array
            if (response.status === 400 && error.errors && Array.isArray(error.errors)) {
              errorMessage = error.errors.join(', ')
            } else {
              errorMessage = error.error || error.message || errorMessage
            }
          } else {
            // If not JSON, use response text or status
            const text = await response.text()
            errorMessage = text || response.statusText || errorMessage
          }
        } catch (e) {
          // If parsing fails, use the response status text
          errorMessage = response.statusText || errorMessage
        }
        // Include status code in error message
        errorMessage = `${response.status}: ${errorMessage}`
        throw new Error(errorMessage)
      }

      // Handle different response types
      if (response.status === 202) {
        // 202 Accepted - may not have response body
        const contentType = response.headers.get('content-type')
        if (contentType && contentType.includes('application/json')) {
          try {
            return await response.json()
          } catch (jsonError) {
            return { success: true, message: 'Request accepted' }
          }
        }
        return { success: true, message: 'Request accepted' }
      }

      // Handle successful responses
      const contentType = response.headers.get('content-type')
      
      // Read response as text first
      const textResponse = await response.text()
      
      // Handle empty responses
      if (!textResponse || textResponse.trim() === '') {
        return { success: true, message: 'Request completed successfully' }
      }
      
      if (contentType && contentType.includes('application/json')) {
        try {
          const jsonResponse = JSON.parse(textResponse)
          return jsonResponse
        } catch (jsonError) {
          // If JSON parsing fails but response was successful, return success
          return { success: true, message: 'Request completed successfully' }
        }
      } else {
        // Return text for non-JSON responses
        return textResponse
      }
    } catch (error) {
      throw error
    }
  }
}

export const api = new ApiClient()