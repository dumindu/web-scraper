// Keywords operations
import { api } from '/model/api.js'

export const keywords = {
  // GET /keywords - Get the list of keywords uploaded by current user
  // Returns: KeywordDTO[]
  // Requires: BearerToken authentication
  async getAll() {
    console.log('Fetching all keywords...')
    return await api.get('/keywords')
  },

  // GET /keywords/{id} - Get the result of a keyword uploaded by current user
  // Parameters: id (string) - Keyword ID in uuid format
  // Returns: KeywordDTO
  // Requires: BearerToken authentication
  async getById(id) {
    console.log('Fetching keyword by ID:', id)
    return await api.get(`/keywords/${id}`)
  },

  // POST /keywords - Upload the keywords CSV file to scrape on web
  // Content-Type: multipart/form-data or text/csv
  // Returns: 202 Accepted
  // Requires: BearerToken authentication
  async uploadCsv(file) {
    console.log('Uploading CSV file:', file.name)
    const formData = new FormData()
    formData.append('file', file)
    
    return await api.postFormData('/keywords', formData)
  },

  // Transform KeywordDTO for display
  // KeywordDTO structure from API:
  // {
  //   id: integer,
  //   keyword: string,
  //   status: string,
  //   adCount: integer,
  //   linkCount: integer,
  //   errorMessage: string,
  //   htmlContent: string
  // }
  formatKeyword(keyword) {
    return {
      id: keyword.id,
      keyword: keyword.keyword,
      status: keyword.status,
      adCount: keyword.adCount || 0,
      linkCount: keyword.linkCount || 0,
      hasError: !!keyword.errorMessage,
      errorMessage: keyword.errorMessage,
      htmlContent: keyword.htmlContent
    }
  },

  // Get status badge class based on keyword status
  getStatusClass(status) {
    switch (status?.toLowerCase()) {
      case 'completed': return 'status-success'
      case 'processing': return 'status-warning'
      case 'failed': return 'status-error'
      default: return 'status-pending'
    }
  },

  // Business logic functions
  async loadKeywords() {
    return await this.getAll()
  },

  filterKeywords(keywordsList, statusFilter, searchFilter) {
    let filtered = keywordsList
    
    if (statusFilter) {
      filtered = filtered.filter(keyword => keyword.status === statusFilter)
    }
    
    if (searchFilter) {
      filtered = filtered.filter(keyword => 
        keyword.keyword.toLowerCase().includes(searchFilter.toLowerCase())
      )
    }
    
    return filtered
  },

  async handleUpload(file) {
    if (!file) {
      throw new Error('Please select a file to upload')
    }
    
    // Validate file type
    if (file.type !== 'text/csv' && !file.name.endsWith('.csv')) {
      throw new Error('Please select a CSV file')
    }
    
    // Validate file size
    if (file.size === 0) {
      throw new Error('Selected file is empty')
    }
    
    if (file.size > 10 * 1024 * 1024) {
      throw new Error('File size must be less than 10MB')
    }
    
    return await this.uploadCsv(file)
  },

  formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }
}