// Dashboard controller functions
import { keywords } from '/model/index.js'

export const dashboardController = {
    async loadDashboardData() {
        try {
            const keywordsList = await keywords.loadKeywords()
            
            // Update stats
            document.getElementById('total-keywords').textContent = keywordsList.length
            document.getElementById('completed-keywords').textContent = 
                keywordsList.filter(k => k.status === 'completed').length
            document.getElementById('processing-keywords').textContent = 
                keywordsList.filter(k => k.status === 'processing').length
            document.getElementById('failed-keywords').textContent = 
                keywordsList.filter(k => k.status === 'failed').length
            
            // Show recent keywords
            const recentKeywords = keywordsList.slice(0, 5)
            const listEl = document.getElementById('recent-keywords-list')
            
            if (recentKeywords.length === 0) {
                listEl.innerHTML = '<p>No keywords uploaded yet. <a href="/upload" onclick="window.navigateTo(\'/upload\'); return false;">Upload some keywords</a> to get started.</p>'
            } else {
                listEl.innerHTML = recentKeywords.map(keyword => {
                    const formatted = keywords.formatKeyword(keyword)
                    return `
                        <div class="keyword-item">
                            <div class="keyword-info">
                                <span class="keyword-text">${formatted.keyword}</span>
                                <span class="keyword-status ${keywords.getStatusClass(formatted.status)}">${formatted.status}</span>
                            </div>
                            <div class="keyword-stats">
                                <span>Links: ${formatted.linkCount}</span>
                                <span>Ads: ${formatted.adCount}</span>
                            </div>
                        </div>
                    `
                }).join('')
            }
        } catch (error) {
            document.getElementById('recent-keywords-list').innerHTML = 
                '<p class="error">Failed to load keywords</p>'
        }
    }
}