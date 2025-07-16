import { auth } from '/model/index.js'
import { dashboardController } from '/controllers/dashboard.js'
import { keywordsController } from '/controllers/keywords.js'
import { uploadController } from '/controllers/upload.js'

export class Router {
    constructor() {
        this.routes = new Map()
        this.currentPath = '/'
    }

    addRoute(path, handler) {
        this.routes.set(path, handler)
    }

    async navigate(path) {
        this.currentPath = path
        window.history.pushState({}, '', path)
        await this.handleRoute(path)
    }

    async handleRoute(path) {
        if (path === '/users/activate') {
            const result = this.showActivateForm()
            if (result && result.email && result.token) {
                await this.autoActivate(result.email, result.token)
            }
            return
        }
        
        // Handle keyword detail routes (e.g., /keywords/123)
        if (path.startsWith('/keywords/')) {
            const keywordId = path.split('/')[2]
            if (keywordId) {
                await this.showKeywordDetail(keywordId)
                return
            }
        }
        
        // Strip query parameters for route matching
        const routePath = path.split('?')[0]
        
        const handler = this.routes.get(routePath) || this.routes.get('*')
        if (handler) {
            await handler()
        }
    }

    init() {
        window.addEventListener('popstate', async () => {
            await this.handleRoute(window.location.pathname)
        })

        this.handleInitialLoad()
    }

    async handleInitialLoad() {
        const path = window.location.pathname
        const isAuthenticated = auth.isAuthenticated()
        
        if (path === '/users/activate') {
            this.hideLoading()
            const result = this.showActivateForm()
            if (result && result.email && result.token) {
                this.autoActivate(result.email, result.token)
            }
            return
        }
        
        this.hideLoading()
        
        // Disable automatic activation check to prevent interference with login form
        // if (isAuthenticated && !await auth.checkActivationStatus()) {
        //     this.navigate('/activate')
        //     return
        // }

        if (isAuthenticated && ['/login', '/signup'].includes(path)) {
            this.navigate('/dashboard')
            return
        }

        if (!isAuthenticated && !['/login', '/signup', '/activate'].includes(path)) {
            this.showLoginForm()
            return
        }

        this.handleRoute(path)
    }

    hideLoading() {
        const loadingScreen = document.getElementById('loading')
        if (loadingScreen) {
            loadingScreen.style.display = 'none'
        }
    }

    async showLoginForm() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Load login screen directly
        await this.loadAuthScreen('login')
        
        appContainer.style.display = 'flex'
        
        // Add event listener for login form after it loads
        setTimeout(() => {
            const loginBtn = document.querySelector('#login-form button')
            if (loginBtn) {
                loginBtn.addEventListener('click', (event) => {
                    event.preventDefault()
                    window.handleSignIn(event)
                })
            }
        }, 100)
    }

    async showSignupForm() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Load signup screen directly
        await this.loadAuthScreen('signup')
        
        appContainer.style.display = 'flex'
    }

    showActivateForm() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Load activate screen directly
        this.loadAuthScreen('activate')
        
        appContainer.style.display = 'flex'
        
        // Handle email display and auto-activation
        const urlParams = new URLSearchParams(window.location.search)
        const emailParam = urlParams.get('email')
        const tokenParam = urlParams.get('token')
        
        if (emailParam && tokenParam) {
            return { email: emailParam, token: tokenParam }
        }
        
        // Get email from URL parameter or localStorage
        const email = emailParam || auth.getPendingActivationEmail()
        
        // Wait for the screen to load then populate fields
        setTimeout(() => {
            const emailDisplay = document.getElementById('email-display')
            const emailHidden = document.getElementById('activation-email')
            
            if (emailDisplay && email) {
                emailDisplay.textContent = email
            }
            
            if (emailHidden && email) {
                emailHidden.value = email
            }
        }, 100)
    }

    async showDashboard() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Initialize app container with header
        await this.initializeAppContainer()
        
        // Load dashboard screen from .dhtml file
        await this.loadScreen('dashboard', 'screen-container')
        
        appContainer.style.display = 'flex'
        
        dashboardController.loadDashboardData()
    }

    async showKeywords() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Initialize app container with header
        await this.initializeAppContainer()
        
        // Load keywords screen from .dhtml file
        await this.loadScreen('keywords', 'screen-container')
        
        appContainer.style.display = 'flex'
        
        keywordsController.loadKeywords()
    }

    async showUpload() {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Initialize app container with header
        await this.initializeAppContainer()
        
        // Load upload screen from .dhtml file
        await this.loadScreen('upload', 'screen-container')
        
        appContainer.style.display = 'flex'
        
        uploadController.resetUploadForm()
        uploadController.setupDragDrop()
    }

    async showKeywordDetail(keywordId) {
        this.hideLoading()
        
        const appContainer = document.getElementById('app-container')
        
        // Initialize app container with header
        await this.initializeAppContainer()
        
        // Load keyword detail screen from .dhtml file
        await this.loadScreen('keyword-detail', 'screen-container')
        
        appContainer.style.display = 'flex'
        
        this.loadKeywordDetail(keywordId)
    }

    async loadKeywordDetail(keywordId) {
        try {
            // This would need to be implemented in the keywords model
            const keyword = await keywordsController.getKeywordById(keywordId)
            const formatted = keywordsController.formatKeyword(keyword)
            
            // Update header
            document.getElementById('keyword-title').textContent = `"${formatted.keyword}"`
            
            // Update status
            const statusBadge = document.getElementById('status-badge')
            statusBadge.textContent = formatted.status
            statusBadge.className = `status-badge ${keywordsController.getStatusClass(formatted.status)}`
            
            // Update stats
            document.getElementById('link-count').textContent = formatted.linkCount
            document.getElementById('ad-count').textContent = formatted.adCount
            
            // Show error if exists
            if (formatted.hasError) {
                document.getElementById('error-section').style.display = 'block'
                document.getElementById('error-message').textContent = formatted.errorMessage
            }
            
            // Show content
            if (formatted.htmlContent) {
                document.getElementById('html-preview').innerHTML = formatted.htmlContent
            } else {
                document.getElementById('html-preview').innerHTML = '<p>No content available</p>'
            }
            
            // Show raw data
            document.getElementById('raw-data').textContent = JSON.stringify(keyword, null, 2)
            
        } catch (error) {
            document.getElementById('html-preview').innerHTML = '<p class="error">Failed to load keyword details</p>'
        }
    }


    async loadScreen(screenName, containerId) {
        const container = document.getElementById(containerId)
        
        try {
            const component = await import(`/view/screens/${screenName}.js`)
            const screenContent = component.default.tmpl
            container.innerHTML = screenContent
            return true
        } catch (error) {
            // Redirect to appropriate route based on screen
            if (screenName === 'dashboard') {
                this.navigate('/keywords')
            } else if (screenName === 'keywords') {
                this.navigate('/dashboard')
            } else if (screenName === 'upload') {
                this.navigate('/dashboard')
            } else {
                this.navigate('/login')
            }
            return false
        }
    }

    async loadAuthScreen(screenName) {
        const appContainer = document.getElementById('app-container')
        
        try {
            const component = await import(`/view/screens/${screenName}.js`)
            const screenContent = component.default.tmpl
            appContainer.innerHTML = screenContent
            return true
        } catch (error) {
            // Redirect to login if auth screen fails to load
            if (screenName !== 'login') {
                this.navigate('/login')
            } else {
                // If login screen fails, show error and redirect to signup
                alert('Login screen failed to load. Redirecting to signup.')
                this.navigate('/signup')
            }
            return false
        }
    }

    async initializeAppContainer() {
        const appContainer = document.getElementById('app-container')
        
        try {
            const component = await import(`/view/layout/app.js`)
            const appLayout = component.default.tmpl
            appContainer.innerHTML = appLayout
            
            // Add screen container inside the main slot
            const main = appContainer.querySelector('main')
            if (main) {
                main.innerHTML = '<div class="main-content"><div id="screen-container"></div></div>'
            }
        } catch (error) {
            // Error loading app layout
        }
    }


    async autoActivate(email, token) {
        try {
            await auth.autoActivate(email, token)
            alert('Account activated successfully!')
            this.navigate('/login')
        } catch (error) {
            alert('Activation failed: ' + error.message)
            this.navigate('/activate')
        }
    }
}

export const router = new Router()

router.addRoute('/login', () => router.showLoginForm())
router.addRoute('/signup', () => router.showSignupForm())
router.addRoute('/activate', () => router.showActivateForm())
router.addRoute('/users/activate', () => router.showActivateForm())
router.addRoute('/dashboard', async () => await router.showDashboard())
router.addRoute('/keywords', async () => await router.showKeywords())
router.addRoute('/upload', async () => await router.showUpload())
router.addRoute('*', () => router.showLoginForm())