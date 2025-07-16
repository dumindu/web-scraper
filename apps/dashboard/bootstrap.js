console.log('Bootstrap.js loading...')

if (typeof window === 'undefined') {
    console.error('Not in browser environment')
}

import { auth } from '/model/index.js'
import { authController } from '/controllers/auth.js'
import { keywordsController } from '/controllers/keywords.js'
import { uploadController } from '/controllers/upload.js'
import { router } from '/controllers/router.js'
import { colorPreferenceController } from '/controllers/color-preference.js'

// Global navigation function
window.navigateTo = (path) => {
    router.navigate(path)
}

// Global logout handler
window.handleLogout = () => {
    auth.logout()
    router.navigate('/login')
}

// Global form handlers
window.handleSignIn = (event) => {
    console.log('Global handleSignIn called')
    event.preventDefault()
    return authController.handleSignIn(event)
}
window.handleSignUp = (event) => authController.handleSignUp(event)
window.handleActivate = (event) => authController.handleActivate(event)
window.resendCode = () => authController.resendCode()
window.showLoginForm = () => router.navigate('/login')
window.showSignupForm = () => router.navigate('/signup')

// Make controller functions globally available
window.filterKeywords = () => keywordsController.filterKeywords()
window.viewKeywordDetails = (id) => keywordsController.viewKeywordDetails(id)
window.handleFileSelect = (event) => uploadController.handleFileSelect(event)
window.clearFile = () => uploadController.clearFile()
window.handleUpload = (event) => uploadController.handleUpload(event)

// Keyword detail functions
window.showTab = (tabName) => {
    // Remove active class from all tabs
    document.querySelectorAll('.tab-button').forEach(button => {
        button.classList.remove('active')
    })
    document.querySelectorAll('.tab-panel').forEach(panel => {
        panel.classList.remove('active')
    })
    
    // Add active class to selected tab
    event.target.classList.add('active')
    document.getElementById(tabName + '-tab').classList.add('active')
}

window.goBack = () => {
    router.navigate('/keywords')
}

// Initialize router
router.init()
colorPreferenceController.handleColorPreference()
