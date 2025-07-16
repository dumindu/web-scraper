// Keywords controller functions
import { keywords } from '/model/index.js'

export const keywordsController = {
    allKeywords: [],

    async loadKeywords() {
        try {
            this.allKeywords = await keywords.loadKeywords()
            this.displayKeywords(this.allKeywords)
        } catch (error) {
            document.getElementById('keywords-list').innerHTML = 
                '<div class="error">Failed to load keywords</div>'
        }
    },

    displayKeywords(keywordsList) {
        const listEl = document.getElementById('keywords-list')
        
        if (keywordsList.length === 0) {
            listEl.innerHTML = '<div class="empty">No keywords found. <a href="/upload" onclick="window.navigateTo(\'/upload\'); return false;">Upload some keywords</a> to get started.</div>'
            return
        }
        
        listEl.innerHTML = keywordsList.map(keyword => {
            const formatted = keywords.formatKeyword(keyword)
            return `
                <div class="table-row">
                    <div class="table-cell">
                        <span class="keyword-text">${formatted.keyword}</span>
                        ${formatted.hasError ? '<span class="error-indicator" title="' + formatted.errorMessage + '">!</span>' : ''}
                    </div>
                    <div class="table-cell">
                        <span class="status-badge ${keywords.getStatusClass(formatted.status)}">${formatted.status}</span>
                    </div>
                    <div class="table-cell">${formatted.linkCount}</div>
                    <div class="table-cell">${formatted.adCount}</div>
                    <div class="table-cell">
                        <button class="btn btn-sm btn-secondary" onclick="viewKeywordDetails('${formatted.id}')">View Details</button>
                    </div>
                </div>
            `
        }).join('')
    },

    filterKeywords() {
        const statusFilter = document.getElementById('status-filter').value
        const searchFilter = document.getElementById('search-filter').value
        
        const filtered = keywords.filterKeywords(this.allKeywords, statusFilter, searchFilter)
        this.displayKeywords(filtered)
    },

    viewKeywordDetails(id) {
        window.navigateTo(`/keywords/${id}`)
    },

    async getKeywordById(id) {
        // This would need to be implemented in the keywords model
        return await keywords.getById(id)
    },

    formatKeyword(keyword) {
        return keywords.formatKeyword(keyword)
    },

    getStatusClass(status) {
        return keywords.getStatusClass(status)
    }
}