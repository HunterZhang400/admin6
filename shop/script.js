// 加载数据
let appData = {};

// 从data.json加载数据
async function loadData() {
    try {
        const response = await fetch('data.json');
        if (response.ok) {
            appData = await response.json();
        } else {
            throw new Error('Failed to load data.json');
        }
    } catch (error) {
        console.error('加载数据失败:', error);
        // 如果加载失败，尝试使用内联数据
        if (window.appData) {
            appData = window.appData;
        } else {
            // 使用默认空数据
            appData = {
                articles: [],
                categories: [],
                rankings: {},
                recommends: []
            };
        }
    }
    initApp();
}

// 页面加载完成后加载数据
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadData);
} else {
    loadData();
}

// 初始化应用
function initApp() {
    renderArticles();
    renderCategoryStats();
    renderRankings('brick');
    renderRecommends();
    initEventListeners();
}

// 渲染文章列表
function renderArticles(articles = appData.articles || []) {
    const articlesList = document.getElementById('articlesList');
    if (!articlesList) return;

    articlesList.innerHTML = articles.map(article => `
        <div class="article-item">
            <img src="${article.thumb}" alt="${article.title}" class="article-thumb" onerror="this.src='https://via.placeholder.com/280x160?text=文章'">
            <div class="article-content">
                <div class="article-date">${article.updateTime}</div>
                <a href="#" class="article-title">${article.title}</a>
                <p class="article-excerpt">${article.excerpt}</p>
                <div class="article-footer">
                    <div class="article-tags">
                        ${article.tags.slice(0, 2).map(tag => `<span class="article-tag">${tag}</span>`).join('')}
                    </div>
                    <div class="article-price">${article.price}</div>
                </div>
                <button class="article-view-btn">点击查看</button>
            </div>
        </div>
    `).join('');

    renderPagination();
}

// 渲染分页
function renderPagination(currentPage = 1, totalPages = 10) {
    const pagination = document.getElementById('pagination');
    if (!pagination) return;

    let paginationHTML = '';
    
    // 上一页
    if (currentPage > 1) {
        paginationHTML += `<a href="#" class="page-btn" data-page="${currentPage - 1}">上一页</a>`;
    }

    // 页码
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);

    if (startPage > 1) {
        paginationHTML += `<a href="#" class="page-btn" data-page="1">1</a>`;
        if (startPage > 2) {
            paginationHTML += `<span class="page-btn">...</span>`;
        }
    }

    for (let i = startPage; i <= endPage; i++) {
        paginationHTML += `<a href="#" class="page-btn ${i === currentPage ? 'active' : ''}" data-page="${i}">${i}</a>`;
    }

    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            paginationHTML += `<span class="page-btn">...</span>`;
        }
        paginationHTML += `<a href="#" class="page-btn" data-page="${totalPages}">${totalPages}</a>`;
    }

    // 下一页
    if (currentPage < totalPages) {
        paginationHTML += `<a href="#" class="page-btn" data-page="${currentPage + 1}">下一页</a>`;
    }

    pagination.innerHTML = paginationHTML;

    // 绑定分页事件
    pagination.querySelectorAll('.page-btn[data-page]').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.preventDefault();
            const page = parseInt(btn.dataset.page);
            // 这里可以调用API加载对应页面的数据
            console.log('切换到第', page, '页');
        });
    });
}

// 渲染分类统计
function renderCategoryStats() {
    const categoryStats = document.getElementById('categoryStats');
    if (!categoryStats || !appData.categories) return;

    categoryStats.innerHTML = appData.categories.map(cat => `
        <div class="category-stat-item">
            <a href="#">${cat.name}</a>
            <span>${cat.count}篇 | ${cat.views}</span>
        </div>
    `).join('');
}

// 渲染榜单
function renderRankings(type = 'brick') {
    const rankingList = document.getElementById('rankingList');
    if (!rankingList || !appData.rankings) return;

    const rankings = appData.rankings[type] || [];
    rankingList.innerHTML = rankings.map((item, index) => `
        <div class="ranking-item">
            <span style="color: var(--primary-color); font-weight: 600; min-width: 20px;">${index + 1}</span>
            <a href="#">${item.title}</a>
            <span class="ranking-heat">${item.heat}热度值</span>
        </div>
    `).join('');
}

// 渲染推荐内容
function renderRecommends() {
    const recommendList = document.getElementById('recommendList');
    if (!recommendList || !appData.recommends) return;

    recommendList.innerHTML = appData.recommends.map(item => `
        <div class="recommend-item">
            <a href="#">${item.title}</a>
        </div>
    `).join('');
}

// 初始化事件监听
function initEventListeners() {
    // 分类标签切换
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', (e) => {
            e.preventDefault();
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            
            const category = tab.dataset.category;
            filterArticles(category);
        });
    });

    // 榜单切换
    const rankTabs = document.querySelectorAll('.rank-tab');
    rankTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            rankTabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            
            const rankType = tab.dataset.rank;
            renderRankings(rankType);
        });
    });

    // 搜索功能
    const searchBtn = document.querySelector('.search-btn');
    const searchInput = document.getElementById('searchInput');
    
    if (searchBtn && searchInput) {
        searchBtn.addEventListener('click', handleSearch);
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                handleSearch();
            }
        });
    }

    // 排序功能（链接形式）
    const sortLinks = document.querySelectorAll('.sort-link');
    sortLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            sortLinks.forEach(l => l.classList.remove('active'));
            link.classList.add('active');
            const sortType = link.dataset.sort;
            sortArticles(sortType);
        });
    });

    // Hero搜索功能
    const heroSearchBtn = document.querySelector('.hero-search-btn');
    const heroSearchInput = document.getElementById('heroSearchInput');
    
    if (heroSearchBtn && heroSearchInput) {
        heroSearchBtn.addEventListener('click', handleHeroSearch);
        heroSearchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                handleHeroSearch();
            }
        });
    }

    // 热门标签点击
    document.querySelectorAll('.hot-tag').forEach(tag => {
        tag.addEventListener('click', (e) => {
            e.preventDefault();
            const keyword = tag.textContent.trim();
            if (heroSearchInput) {
                heroSearchInput.value = keyword;
                handleHeroSearch();
            }
        });
    });

    // 返回顶部
    const backToTop = document.getElementById('backToTop');
    if (backToTop) {
        backToTop.addEventListener('click', (e) => {
            e.preventDefault();
            window.scrollTo({ top: 0, behavior: 'smooth' });
        });

        // 滚动显示/隐藏返回顶部按钮
        window.addEventListener('scroll', () => {
            if (window.scrollY > 300) {
                backToTop.style.display = 'flex';
            } else {
                backToTop.style.display = 'none';
            }
        });
    }

    // 主题切换
    const themeToggle = document.getElementById('themeToggle');
    if (themeToggle) {
        themeToggle.addEventListener('click', (e) => {
            e.preventDefault();
            document.body.classList.toggle('dark-theme');
            const isDark = document.body.classList.contains('dark-theme');
            themeToggle.textContent = isDark ? '☀' : '🌓';
        });
    }
}

// Hero搜索处理
function handleHeroSearch() {
    const heroSearchInput = document.getElementById('heroSearchInput');
    if (!heroSearchInput || !appData.articles) return;
    
    const keyword = heroSearchInput.value.trim().toLowerCase();
    
    if (!keyword) {
        renderArticles();
        return;
    }
    
    const filtered = appData.articles.filter(article => 
        article.title.toLowerCase().includes(keyword) ||
        article.excerpt.toLowerCase().includes(keyword) ||
        article.tags.some(tag => tag.toLowerCase().includes(keyword))
    );
    
    renderArticles(filtered);
}

// 筛选文章
function filterArticles(category) {
    if (!appData.articles) return;
    
    let filtered = appData.articles;
    
    if (category && category !== 'all') {
        const categoryMap = {
            'free': '免费资源',
            'brick': '搬砖项目',
            'earn': '网赚项目',
            'script': '挂机脚本',
            'media': '自媒体类',
            'ecommerce': '电商运营',
            'other': '其他分类'
        };
        
        const categoryName = categoryMap[category];
        if (categoryName) {
            filtered = appData.articles.filter(article => 
                article.category === categoryName || 
                article.tags.some(tag => tag.includes(categoryName))
            );
        }
    }
    
    renderArticles(filtered);
}

// 搜索文章
function handleSearch() {
    const searchInput = document.getElementById('searchInput');
    if (!searchInput || !appData.articles) return;
    
    const keyword = searchInput.value.trim().toLowerCase();
    
    if (!keyword) {
        renderArticles();
        return;
    }
    
    const filtered = appData.articles.filter(article => 
        article.title.toLowerCase().includes(keyword) ||
        article.excerpt.toLowerCase().includes(keyword) ||
        article.tags.some(tag => tag.toLowerCase().includes(keyword))
    );
    
    renderArticles(filtered);
}

// 排序文章
function sortArticles(sortType) {
    if (!appData.articles) return;
    
    const sorted = [...appData.articles];
    
    switch(sortType) {
        case 'update':
            sorted.sort((a, b) => new Date(b.updateTime) - new Date(a.updateTime));
            break;
        case 'publish':
            sorted.sort((a, b) => new Date(b.updateTime) - new Date(a.updateTime));
            break;
        case 'view':
            sorted.sort((a, b) => b.views - a.views);
            break;
        case 'like':
            sorted.sort((a, b) => b.likes - a.likes);
            break;
        case 'comment':
            sorted.sort((a, b) => b.comments - a.comments);
            break;
        case 'random':
            sorted.sort(() => Math.random() - 0.5);
            break;
    }
    
    renderArticles(sorted);
}

// 响应式菜单切换（移动端）
function initMobileMenu() {
    const nav = document.querySelector('.main-nav');
    if (window.innerWidth <= 768) {
        // 移动端菜单处理
    }
}

window.addEventListener('resize', initMobileMenu);
initMobileMenu();
