// Upload controller functions
import { keywords } from '/model/index.js'

export const uploadController = {
    selectedFile: null,

    validateFile(file) {
        const errorEl = document.getElementById('upload-error-message')
        
        // Validate file type
        if (file.type !== 'text/csv' && !file.name.endsWith('.csv')) {
            errorEl.textContent = 'Please select a CSV file'
            errorEl.style.display = 'block'
            this.resetUploadForm()
            return false
        }
        
        // Validate file size (minimum 1 byte, maximum 10MB)
        if (file.size === 0) {
            errorEl.textContent = 'Selected file is empty'
            errorEl.style.display = 'block'
            this.resetUploadForm()
            return false
        }
        
        if (file.size > 10 * 1024 * 1024) {
            errorEl.textContent = 'File size must be less than 10MB'
            errorEl.style.display = 'block'
            this.resetUploadForm()
            return false
        }
        
        return true
    },

    handleFileSelect(event) {
        const file = event.target.files[0]
        const errorEl = document.getElementById('upload-error-message')
        
        // Clear previous errors
        errorEl.style.display = 'none'
        
        if (file && this.validateFile(file)) {
            this.selectedFile = file
            this.showFileInfo(file)
            document.getElementById('upload-button').disabled = false
        }
    },

    showFileInfo(file) {
        const fileInfo = document.getElementById('file-info')
        const fileName = document.getElementById('file-name')
        const fileSize = document.getElementById('file-size')
        
        fileName.textContent = file.name
        fileSize.textContent = keywords.formatFileSize(file.size)
        fileInfo.style.display = 'flex'
    },

    clearFile() {
        this.resetUploadForm()
    },

    async handleUpload(event) {
        event.preventDefault()
        
        const errorEl = document.getElementById('upload-error-message')
        const uploadButton = document.getElementById('upload-button')
        
        // Clear previous errors
        errorEl.style.display = 'none'
        
        // Disable upload button and show loading state
        uploadButton.disabled = true
        uploadButton.textContent = 'Uploading...'
        
        try {
            const response = await keywords.handleUpload(this.selectedFile)
            
            // Show success message
            const successMessage = document.createElement('div')
            successMessage.className = 'form-success'
            successMessage.style.display = 'block'
            successMessage.style.color = 'var(--success, #22c55e)'
            successMessage.style.padding = '0.75rem'
            successMessage.style.marginBottom = '1rem'
            successMessage.style.border = '1px solid var(--success, #22c55e)'
            successMessage.style.borderRadius = '4px'
            successMessage.style.backgroundColor = 'var(--success-bg, #f0fdf4)'
            successMessage.textContent = 'Keywords uploaded successfully! Processing has started.'
            
            // Insert success message before the error message element
            errorEl.parentNode.insertBefore(successMessage, errorEl)
            
            // Reset form after a short delay
            setTimeout(() => {
                this.resetUploadForm()
                window.navigateTo('/keywords')
            }, 2000)
            
        } catch (error) {
            // Show error message
            errorEl.textContent = error.message || 'Upload failed. Please try again.'
            errorEl.style.display = 'block'
            
            // Re-enable upload button
            uploadButton.disabled = false
            uploadButton.textContent = 'Upload Keywords'
        }
    },

    resetUploadForm() {
        this.selectedFile = null
        document.getElementById('file-input').value = ''
        document.getElementById('file-info').style.display = 'none'
        document.getElementById('upload-button').disabled = true
        document.getElementById('upload-error-message').style.display = 'none'
        
        // Remove any success messages
        const successMessages = document.querySelectorAll('.form-success')
        successMessages.forEach(msg => msg.remove())
        
        // Reset upload button text
        document.getElementById('upload-button').textContent = 'Upload Keywords'
    },

    setupDragDrop() {
        const uploadArea = document.getElementById('upload-area')
        if (!uploadArea) return
        
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault()
            e.stopPropagation()
            uploadArea.classList.add('drag-over')
        })
        
        uploadArea.addEventListener('dragleave', (e) => {
            e.preventDefault()
            e.stopPropagation()
            uploadArea.classList.remove('drag-over')
        })
        
        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault()
            e.stopPropagation()
            uploadArea.classList.remove('drag-over')
            
            const files = e.dataTransfer.files
            if (files.length > 0) {
                const file = files[0]
                const errorEl = document.getElementById('upload-error-message')
                
                // Clear previous errors
                errorEl.style.display = 'none'
                
                if (this.validateFile(file)) {
                    this.selectedFile = file
                    this.showFileInfo(file)
                    document.getElementById('upload-button').disabled = false
                }
            }
        })
        
        uploadArea.addEventListener('click', () => {
            document.getElementById('file-input').click()
        })
    }
}