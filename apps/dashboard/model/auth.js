// Authentication operations
import { api } from '/model/api.js'

export const auth = {
  async signUp(email, password, confirmPassword) {
    const response = await api.post('/users/sign-up', {
      email,
      password,
      confirmPassword
    })
    
    // Store email for activation process on successful signup
    if (response && (response.success || response.message)) {
      localStorage.setItem('pending_activation_email', email)
    }
    
    return response
  },

  async signIn(email, password) {
    try {
      const response = await api.post('/users/sign-in', {
        email,
        password
      })
      
      console.log('Login response:', response)
      
      if (response.access && response.refresh) {
        api.setToken(response.access)
        localStorage.setItem('refresh_token', response.refresh)
        localStorage.setItem('user_email', email)
        localStorage.removeItem('pending_activation_email')
        
        console.log('Tokens stored successfully')
        console.log('Access token:', response.access.substring(0, 50) + '...')
        console.log('Refresh token:', response.refresh.substring(0, 50) + '...')
      }
      
      return response
    } catch (error) {
      console.log('Auth.signIn caught error:', error.message)
      console.log('Error includes 403:', error.message.includes('403'))
      console.log('Error includes pending activation:', error.message.includes('pending activation'))
      
      // Only store email for activation if error specifically mentions pending activation
      if (error.message === 'pending activation' || error.message.includes('pending activation')) {
        console.log('Storing email for activation:', email)
        // Store email for activation process
        localStorage.setItem('pending_activation_email', email)
      }
      console.log('Rethrowing error to controller')
      throw error
    }
  },

  async activate(email, token) {
    const response = await api.post(`/users/activate?email=${encodeURIComponent(email)}&token=${encodeURIComponent(token)}`)
    
    if (response) {
      localStorage.removeItem('pending_activation_email')
    }
    
    return response
  },

  logout() {
    api.clearToken()
    localStorage.removeItem('user_email')
    localStorage.removeItem('pending_activation_email')
  },

  isAuthenticated() {
    return !!localStorage.getItem('access_token')
  },

  getToken() {
    return localStorage.getItem('access_token')
  },

  getUserEmail() {
    return localStorage.getItem('user_email')
  },

  getPendingActivationEmail() {
    return localStorage.getItem('pending_activation_email')
  },

  async checkActivationStatus() {
    // If user has valid token, they should be activated
    // If they get 403 responses, they need activation
    const token = this.getToken()
    if (!token) return false
    
    try {
      // Try to make a test request to check if user is activated
      await api.get('/keywords')
      return true
    } catch (error) {
      // If we get 403, user needs activation
      if (error.message.includes('403') || error.message.includes('Forbidden')) {
        return false
      }
      return true // Other errors assume user is activated
    }
  },

  needsActivation() {
    return !!this.getPendingActivationEmail()
  },

  // Form handlers
  async handleSignIn(email, password) {
    await this.signIn(email, password)
  },

  async handleSignUp(email, password, confirmPassword) {
    if (password !== confirmPassword) {
      throw new Error('Passwords do not match')
    }
    await this.signUp(email, password, confirmPassword)
  },

  async handleActivate(email, token) {
    const response = await api.post(`/users/activate?email=${encodeURIComponent(email)}&token=${encodeURIComponent(token)}`)
    
    if (response) {
      localStorage.removeItem('pending_activation_email')
    }
    
    return response
  },

  async resendCode() {
    const email = this.getPendingActivationEmail()
    if (!email) {
      throw new Error('No pending activation found')
    }
    
    await this.signUp(email, '', '')
    return { message: 'Activation code resent to your email' }
  },

  async autoActivate(email, token) {
    const response = await api.post(`/users/activate?email=${encodeURIComponent(email)}&token=${encodeURIComponent(token)}`)
    
    if (response) {
      localStorage.removeItem('pending_activation_email')
    }
    
    return response
  }
}