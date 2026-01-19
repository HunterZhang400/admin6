// 六合彩标准规则配置
// 注：生肖对应数字每年会随农历新年变化，此处使用标准基础规则

// 标准生肖对应数字（以某年度为例，实际应随年份更新）
const zodiacNumbers = {
    '鼠': [4, 16, 28, 40],
    '牛': [3, 15, 27, 39],
    '虎': [2, 14, 26, 38],
    '兔': [1, 13, 25, 37, 49],
    '龙': [12, 24, 36, 48],
    '蛇': [11, 23, 35, 47],
    '马': [10, 22, 34, 46],
    '羊': [9, 21, 33, 45],
    '猴': [8, 20, 32, 44],
    '鸡': [7, 19, 31, 43],
    '狗': [6, 18, 30, 42],
    '猪': [5, 17, 29, 41]
};

// 标准五行分类（根据需求文档：基于个位数规则）
// 规则：个位1或2=木，个位3或4=火，个位5=土，个位6=水，个位7或8=金，个位9=金，个位0=水，49号=金（特例）
const elementNumbers = {
    '金': [7, 8, 9, 17, 18, 19, 27, 28, 29, 37, 38, 39, 47, 48, 49], // 个位7,8,9 + 49特例
    '木': [1, 2, 11, 12, 21, 22, 31, 32, 41, 42], // 个位1或2
    '水': [6, 10, 16, 20, 26, 30, 36, 40, 46], // 个位6或0
    '火': [3, 4, 13, 14, 23, 24, 33, 34, 43, 44], // 个位3或4
    '土': [5, 15, 25, 35, 45] // 个位5
};

// 标准波色分类（固定规则）
const waveNumbers = {
    '红波': [1, 2, 7, 8, 12, 13, 18, 19, 23, 24, 29, 30, 34, 35, 40, 45, 46],
    '蓝波': [3, 4, 9, 10, 14, 15, 20, 25, 26, 31, 36, 37, 41, 42, 47, 48],
    '绿波': [5, 6, 11, 16, 17, 21, 22, 27, 28, 32, 33, 38, 39, 43, 44, 49]
};

// 根据数字获取生肖
function getZodiacByNumber(num) {
    for (const [zodiac, numbers] of Object.entries(zodiacNumbers)) {
        if (numbers.includes(num)) {
            return zodiac;
        }
    }
    // 如果不在列表中，使用循环算法
    const zodiacs = ['鼠', '牛', '虎', '兔', '龙', '蛇', '马', '羊', '猴', '鸡', '狗', '猪'];
    return zodiacs[(num - 1) % 12];
}

// 根据数字获取五行（基于个位数规则，符合需求文档）
function getElementByNumber(num) {
    // 特例：49号属金
    if (num === 49) return '金';
    
    const lastDigit = num % 10;
    
    // 根据个位数判断五行
    if (lastDigit === 1 || lastDigit === 2) return '木';
    if (lastDigit === 3 || lastDigit === 4) return '火';
    if (lastDigit === 5) return '土';
    if (lastDigit === 6 || lastDigit === 0) return '水';
    if (lastDigit === 7 || lastDigit === 8 || lastDigit === 9) return '金';
    
    return '土'; // 默认（理论上不会到这里）
}

// 根据数字获取波色
function getWaveByNumber(num) {
    if (waveNumbers['红波'].includes(num)) return 'red';
    if (waveNumbers['蓝波'].includes(num)) return 'blue';
    if (waveNumbers['绿波'].includes(num)) return 'green';
    // 如果不在列表中，使用算法
    const waves = ['red', 'blue', 'green'];
    return waves[(num - 1) % 3];
}

// 获取数字数据的辅助函数
function getNumberData(num) {
    return {
        color: getWaveByNumber(num),
        zodiac: getZodiacByNumber(num),
        element: getElementByNumber(num)
    };
}

// 生肖分类
const zodiacCategories = {
    '家禽': ['鸡', '鸭', '鹅', '狗', '猪'],
    '野兽': ['虎', '龙', '蛇', '马', '羊', '猴', '鼠', '牛', '兔'],
    '前肖': ['鼠', '牛', '虎', '兔', '龙', '蛇'],
    '后肖': ['马', '羊', '猴', '鸡', '狗', '猪'],
    '天肖': ['牛', '兔', '龙', '马', '猴', '猪'],
    '地肖': ['鼠', '虎', '蛇', '羊', '鸡', '狗'],
    '男肖': ['鼠', '虎', '龙', '马', '猴', '狗'],
    '女肖': ['牛', '兔', '蛇', '羊', '鸡', '猪'],
    '阳肖': ['鼠', '虎', '龙', '马', '猴', '狗'],
    '阴肖': ['牛', '兔', '蛇', '羊', '鸡', '猪']
};

// 内围码和外围码（标准规则：1-24为内围，25-49为外围）
const innerCodes = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24];
const outerCodes = [25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49];

// 初始化页面
document.addEventListener('DOMContentLoaded', function() {
    initRowLabels();
    initNumberGrid();
    initCategoryButtons();
});

// 初始化行标签
function initRowLabels() {
    const rowLabelsContainer = document.getElementById('rowLabels');
    for (let i = 1; i <= 10; i++) {
        const label = document.createElement('div');
        label.className = 'row-label';
        label.textContent = String(i).padStart(2, '0');
        rowLabelsContainer.appendChild(label);
    }
}

// 初始化数字网格
function initNumberGrid() {
    const gridContainer = document.getElementById('numberGrid');
    
    for (let i = 1; i <= 49; i++) {
        const numberItem = document.createElement('div');
        numberItem.className = 'number-item';
        numberItem.dataset.number = i;
        
        const data = getNumberData(i);
        const colorClass = data.color || 'red';
        
        numberItem.innerHTML = `
            <div class="number-circle ${colorClass}" data-number="${i}">
                ${String(i).padStart(2, '0')}
            </div>
            <input type="number" 
                   class="bet-input" 
                   id="input-${i}" 
                   data-number="${i}"
                   placeholder="金额" 
                   min="0" 
                   step="1">
        `;
        
        gridContainer.appendChild(numberItem);
    }
}

// 初始化分类按钮
function initCategoryButtons() {
    const buttons = document.querySelectorAll('.category-btn');
    buttons.forEach(button => {
        button.addEventListener('click', function() {
            const category = this.dataset.category;
            const value = this.dataset.value;
            
            // 切换按钮激活状态
            if (this.classList.contains('active')) {
                this.classList.remove('active');
                clearHighlights();
            } else {
                // 移除其他按钮的激活状态（同一类别）
                document.querySelectorAll(`.category-btn[data-category="${category}"]`).forEach(btn => {
                    btn.classList.remove('active');
                });
                this.classList.add('active');
                highlightNumbers(category, value);
            }
        });
    });
}

// 清除所有高亮
function clearHighlights() {
    document.querySelectorAll('.bet-input').forEach(input => {
        input.classList.remove('highlight');
    });
    document.querySelectorAll('.number-circle').forEach(circle => {
        circle.classList.remove('highlight');
    });
}

// 高亮显示匹配的数字
function highlightNumbers(category, value) {
    clearHighlights();
    
    const matchingNumbers = getMatchingNumbers(category, value);
    
    matchingNumbers.forEach(num => {
        const input = document.getElementById(`input-${num}`);
        const circle = document.querySelector(`.number-circle[data-number="${num}"]`);
        if (input) input.classList.add('highlight');
        if (circle) circle.classList.add('highlight');
    });
}

// 获取匹配的数字列表
function getMatchingNumbers(category, value) {
    const matches = [];
    
    for (let i = 1; i <= 49; i++) {
        const data = getNumberData(i);
        const isOdd = i % 2 === 1;
        const isEven = i % 2 === 0;
        const isBig = i >= 25 && i <= 49;
        const isSmall = i >= 1 && i <= 24;
        const head = Math.floor(i / 10);
        const tail = i % 10;
        const tailBig = tail >= 5;
        const tailSmall = tail < 5;
        const sum = getSumOfDigits(i);
        const sumOdd = sum % 2 === 1;
        const sumEven = sum % 2 === 0;
        const sumBig = sum >= 7;
        const sumSmall = sum < 7;
        
        let match = false;
        
        switch (category) {
            case 'element':
                // 五行分类：使用标准配置
                match = elementNumbers[value] && elementNumbers[value].includes(i);
                break;
                
            case 'wave':
                // 波色分类：使用标准配置
                if (value === '红波') match = waveNumbers['红波'].includes(i);
                else if (value === '蓝波') match = waveNumbers['蓝波'].includes(i);
                else if (value === '绿波') match = waveNumbers['绿波'].includes(i);
                break;
                
            case 'parity':
                if (value === '单') match = isOdd;
                else if (value === '双') match = isEven;
                break;
                
            case 'size':
                if (value === '大') match = isBig;
                else if (value === '小') match = isSmall;
                break;
                
            case 'combo':
                match = checkCombo(value, data, isOdd, isEven, isBig, isSmall, head, tail);
                break;
                
            case 'combined':
                if (value === '合单') match = sumOdd;
                else if (value === '合双') match = sumEven;
                else if (value === '合大') match = sumBig;
                else if (value === '合小') match = sumSmall;
                break;
                
            case 'tail':
                match = checkTail(value, tail, tailBig, tailSmall, sum);
                break;
                
            case 'zodiac':
                // 生肖分类：使用标准配置
                match = zodiacNumbers[value] && zodiacNumbers[value].includes(i);
                break;
                
            case 'zodiac-type':
                match = checkZodiacType(value, data.zodiac);
                break;
                
            case 'head':
                const headNum = parseInt(value);
                match = head === headNum;
                break;
                
            case 'sum':
                const sumNum = parseInt(value);
                match = sum === sumNum;
                break;
                
            case 'code':
                if (value === '内围码') match = innerCodes.includes(i);
                else if (value === '外围码') match = outerCodes.includes(i);
                break;
        }
        
        if (match) {
            matches.push(i);
        }
    }
    
    return matches;
}

// 检查组合条件
function checkCombo(value, data, isOdd, isEven, isBig, isSmall, head, tail) {
    switch (value) {
        case '小单': return isSmall && isOdd;
        case '小双': return isSmall && isEven;
        case '大单': return isBig && isOdd;
        case '大双': return isBig && isEven;
        case '红单': return data.color === 'red' && isOdd;
        case '红双': return data.color === 'red' && isEven;
        case '蓝单': return data.color === 'blue' && isOdd;
        case '蓝双': return data.color === 'blue' && isEven;
        case '绿单': return data.color === 'green' && isOdd;
        case '绿双': return data.color === 'green' && isEven;
        case '红大': return data.color === 'red' && isBig;
        case '红小': return data.color === 'red' && isSmall;
        case '蓝大': return data.color === 'blue' && isBig;
        case '蓝小': return data.color === 'blue' && isSmall;
        case '绿大': return data.color === 'green' && isBig;
        case '绿小': return data.color === 'green' && isSmall;
        case '0头单': return head === 0 && isOdd;
        case '0头双': return head === 0 && isEven;
        case '1头单': return head === 1 && isOdd;
        case '1头双': return head === 1 && isEven;
        case '2头单': return head === 2 && isOdd;
        case '2头双': return head === 2 && isEven;
        case '3头单': return head === 3 && isOdd;
        case '3头双': return head === 3 && isEven;
        case '4头单': return head === 4 && isOdd;
        case '4头双': return head === 4 && isEven;
        default: return false;
    }
}

// 检查尾数条件
function checkTail(value, tail, tailBig, tailSmall, sum) {
    switch (value) {
        case '尾大': return tailBig;
        case '尾小': return tailSmall;
        case '合尾大': {
            const sumTail = sum % 10;
            return sumTail >= 5;
        }
        case '合尾小': {
            const sumTail = sum % 10;
            return sumTail < 5;
        }
        case '0尾': return tail === 0;
        case '1尾': return tail === 1;
        case '2尾': return tail === 2;
        case '3尾': return tail === 3;
        case '4尾': return tail === 4;
        case '5尾': return tail === 5;
        case '6尾': return tail === 6;
        case '7尾': return tail === 7;
        case '8尾': return tail === 8;
        case '9尾': return tail === 9;
        default: return false;
    }
}

// 检查生肖类型
function checkZodiacType(type, zodiac) {
    const category = zodiacCategories[type];
    if (!category) return false;
    return category.includes(zodiac);
}

// 获取数字的合数（数字根）- 符合需求文档：反复相加直到个位数
function getSumOfDigits(num) {
    // 数字根算法：反复将各位数字相加，直到得到一个个位数
    let sum = num;
    while (sum >= 10) {
        const str = sum.toString();
        sum = 0;
        for (let i = 0; i < str.length; i++) {
            sum += parseInt(str[i]);
        }
    }
    return sum;
}
