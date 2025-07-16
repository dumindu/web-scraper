// Auth controller functions
import { auth } from '/model/index.js'

export const authController = {
    async handleSignIn(event) {
        event.preventDefault()
        
        const form = document.getElementById('login-form')
        const formData = new FormData(form)
        const email = formData.get('email')
        const password = formData.get('password')
        
        const errorEl = document.getElementById('login-error-message')
        const submitBtn = form.querySelector('button') || document.querySelector('#login-form button')
        
        errorEl.style.display = 'none'
        submitBtn.disabled = true
        submitBtn.textContent = 'Signing in...'
        
        try {
            await auth.signIn(email, password)
            window.navigateTo('/dashboard')
        } catch (error) {
            // Show pending activation error on form instead of redirecting
            if (error.message.includes('pending activation')) {
                localStorage.setItem('pending_activation_email', email)
                
                errorEl.textContent = 'Account needs activation. Please check your email for the activation link.'
                errorEl.style.display = 'block'
                submitBtn.disabled = false
                submitBtn.textContent = 'Sign In'
                return
            }
            
            // For 401 (Unauthorized), 403 (Forbidden), and 409 (Conflict) errors, just show on form
            if (error.message.includes('401') || error.message.includes('403') || error.message.includes('409')) {
                errorEl.textContent = error.message
                errorEl.style.display = 'block'
                submitBtn.disabled = false
                submitBtn.textContent = 'Sign In'
                return
            }
            
            // For all other errors, show on form
            errorEl.textContent = error.message
            errorEl.style.display = 'block'
            submitBtn.disabled = false
            submitBtn.textContent = 'Sign In'
        }
    },

    async handleSignUp(event) {
        event.preventDefault()
        
        const form = event.target
        const formData = new FormData(form)
        const email = formData.get('email')
        const password = formData.get('password')
        const confirmPassword = formData.get('confirmPassword')
        
        const errorEl = document.getElementById('signup-error-message')
        const submitBtn = form.querySelector('button[type="submit"]')
        
        errorEl.style.display = 'none'
        submitBtn.disabled = true
        submitBtn.textContent = 'Creating account...'
        
        try {
            const response = await auth.signUp(email, password, confirmPassword)
            
            // Ensure email is stored for activation
            localStorage.setItem('pending_activation_email', email)
            window.navigateTo(`/activate?email=${encodeURIComponent(email)}`)
        } catch (error) {
            
            // Check if it's a conflict error (user already exists)
            if (error.message && (error.message.includes('409') || error.message.includes('conflict') || error.message.includes('already exists'))) {
                alert('User already exists. Redirecting to login.')
                window.navigateTo('/login')
            } else {
                errorEl.textContent = error.message
                errorEl.style.display = 'block'
                submitBtn.disabled = false
                submitBtn.textContent = 'Sign Up'
            }
        }
    },

    async handleActivate(event) {
        event.preventDefault()
        
        const form = event.target
        const formData = new FormData(form)
        const token = formData.get('token')
        
        const errorEl = document.getElementById('activate-error-message')
        const submitBtn = form.querySelector('button[type="submit"]')
        
        // Get email from localStorage or URL params
        let email = auth.getPendingActivationEmail()
        if (!email) {
            const urlParams = new URLSearchParams(window.location.search)
            email = urlParams.get('email')
        }
        
        if (!email) {
            errorEl.textContent = 'No pending activation found'
            errorEl.style.display = 'block'
            return
        }
        
        errorEl.style.display = 'none'
        submitBtn.disabled = true
        submitBtn.textContent = 'Activating...'
        
        try {
            await auth.activate(email, token)
            alert('Account activated successfully!')
            window.navigateTo('/login')
        } catch (error) {
            errorEl.textContent = error.message
            errorEl.style.display = 'block'
            submitBtn.disabled = false
            submitBtn.textContent = 'Activate Account'
        }
    },

    async resendCode() {
        try {
            const result = await auth.resendCode()
            alert(result.message)
        } catch (error) {
            alert('Failed to resend code: ' + error.message)
        }
    },

    handleLogout() {
        auth.logout()
        window.navigateTo('/login')
    }
}