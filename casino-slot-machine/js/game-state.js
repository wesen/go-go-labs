/**
 * Game State Management
 * Handles player balance, betting, and game state persistence
 */

class GameState {
    constructor() {
        this.balance = 1000;
        this.currentBet = 10;
        this.lastWin = 0;
        this.totalWins = 0;
        this.totalSpins = 0;
        this.gameHistory = [];
        this.maxHistorySize = 100;
        
        this.symbols = ['🍒', '🍋', '🍊', '🔔', '💎', '⭐', '🍇', '🍀'];
        this.payouts = {
            '🍒🍒🍒': 20,
            '🍋🍋🍋': 15,
            '🍊🍊🍊': 10,
            '🔔🔔🔔': 50,
            '💎💎💎': 100,
            '⭐⭐⭐': 75,
            '🍇🍇🍇': 30,
            '🍀🍀🍀': 200,
            // Two symbol matches
            '🍒🍒': 2,
            '🍋🍋': 2,
            '🍊🍊': 2,
            '🔔🔔': 5,
            '💎💎': 10,
            '⭐⭐': 8,
            '🍇🍇': 3,
            '🍀🍀': 15
        };
        
        this.loadGameState();
        this.bindEvents();
    }
    
    bindEvents() {
        // Save game state when page is about to unload
        window.addEventListener('beforeunload', () => {
            this.saveGameState();
        });
        
        // Auto-save every 30 seconds
        setInterval(() => {
            this.saveGameState();
        }, 30000);
    }
    
    // Balance management
    getBalance() {
        return this.balance;
    }
    
    setBalance(amount) {
        this.balance = Math.max(0, amount);
        this.updateDisplay();
        this.saveGameState();
    }
    
    addToBalance(amount) {
        this.balance += amount;
        this.updateDisplay();
        this.saveGameState();
    }
    
    subtractFromBalance(amount) {
        if (this.balance >= amount) {
            this.balance -= amount;
            this.updateDisplay();
            this.saveGameState();
            return true;
        }
        return false;
    }
    
    // Betting management
    getCurrentBet() {
        return this.currentBet;
    }
    
    setCurrentBet(amount) {
        const maxBet = Math.min(100, this.balance);
        this.currentBet = Math.max(1, Math.min(amount, maxBet));
        this.updateDisplay();
        return this.currentBet;
    }
    
    canAffordBet() {
        return this.balance >= this.currentBet;
    }
    
    // Advanced betting features
    getMaxBet() {
        return Math.min(100, this.balance);
    }
    
    getBetMultipliers() {
        return [1, 5, 10, 25, 100];
    }
    
    autoAdjustBet() {
        // Auto-adjust bet if current bet exceeds balance
        if (this.currentBet > this.balance) {
            this.currentBet = Math.min(this.balance, 10);
            this.updateDisplay();
        }
    }
    
    // Risk management
    getRiskLevel() {
        const betToBalanceRatio = this.currentBet / this.balance;
        if (betToBalanceRatio > 0.1) return 'high';
        if (betToBalanceRatio > 0.05) return 'medium';
        return 'low';
    }
    
    // Win management
    getLastWin() {
        return this.lastWin;
    }
    
    setLastWin(amount) {
        this.lastWin = amount;
        if (amount > 0) {
            this.totalWins += amount;
            this.addToBalance(amount);
        }
        this.updateDisplay();
    }
    
    // Game statistics
    getStats() {
        return {
            balance: this.balance,
            totalWins: this.totalWins,
            totalSpins: this.totalSpins,
            winRate: this.totalSpins > 0 ? (this.getWinCount() / this.totalSpins * 100) : 0,
            biggestWin: this.getBiggestWin()
        };
    }
    
    getWinCount() {
        return this.gameHistory.filter(spin => spin.payout > 0).length;
    }
    
    getBiggestWin() {
        return Math.max(0, ...this.gameHistory.map(spin => spin.payout));
    }
    
    // Spin processing
    processSpin(reelResults) {
        if (!this.canAffordBet()) {
            throw new Error('Insufficient balance');
        }
        
        // Deduct bet from balance
        this.subtractFromBalance(this.currentBet);
        this.totalSpins++;
        
        // Calculate payout
        const payout = this.calculatePayout(reelResults);
        this.setLastWin(payout);
        
        // Record in history
        this.addToHistory({
            spin: this.totalSpins,
            bet: this.currentBet,
            result: reelResults,
            payout: payout,
            timestamp: new Date().toISOString()
        });
        
        return {
            payout: payout,
            newBalance: this.balance,
            isWin: payout > 0
        };
    }
    
    calculatePayout(reelResults) {
        if (!reelResults || reelResults.length !== 3) {
            return 0;
        }
        
        const resultString = reelResults.join('');
        
        // Check for three-symbol matches first
        if (this.payouts[resultString]) {
            return this.payouts[resultString] * this.currentBet;
        }
        
        // Check for two-symbol matches
        const first = reelResults[0];
        const second = reelResults[1];
        const third = reelResults[2];
        
        if (first === second) {
            const twoSymbolKey = first + second;
            if (this.payouts[twoSymbolKey]) {
                return this.payouts[twoSymbolKey] * this.currentBet;
            }
        }
        
        if (second === third) {
            const twoSymbolKey = second + third;
            if (this.payouts[twoSymbolKey]) {
                return this.payouts[twoSymbolKey] * this.currentBet;
            }
        }
        
        if (first === third) {
            const twoSymbolKey = first + third;
            if (this.payouts[twoSymbolKey]) {
                return this.payouts[twoSymbolKey] * this.currentBet;
            }
        }
        
        return 0;
    }
    
    // History management
    addToHistory(spinData) {
        this.gameHistory.unshift(spinData);
        if (this.gameHistory.length > this.maxHistorySize) {
            this.gameHistory = this.gameHistory.slice(0, this.maxHistorySize);
        }
    }
    
    getHistory(limit = 10) {
        return this.gameHistory.slice(0, limit);
    }
    
    clearHistory() {
        this.gameHistory = [];
        this.saveGameState();
    }
    
    // Persistence
    saveGameState() {
        const state = {
            balance: this.balance,
            currentBet: this.currentBet,
            lastWin: this.lastWin,
            totalWins: this.totalWins,
            totalSpins: this.totalSpins,
            gameHistory: this.gameHistory,
            lastSaved: new Date().toISOString()
        };
        
        try {
            localStorage.setItem('slotMachineGameState', JSON.stringify(state));
        } catch (error) {
            console.warn('Could not save game state to localStorage:', error);
        }
    }
    
    loadGameState() {
        try {
            const savedState = localStorage.getItem('slotMachineGameState');
            if (savedState) {
                const state = JSON.parse(savedState);
                this.balance = state.balance || 1000;
                this.currentBet = state.currentBet || 10;
                this.lastWin = state.lastWin || 0;
                this.totalWins = state.totalWins || 0;
                this.totalSpins = state.totalSpins || 0;
                this.gameHistory = state.gameHistory || [];
            }
        } catch (error) {
            console.warn('Could not load game state from localStorage:', error);
            this.resetGameState();
        }
    }
    
    resetGameState() {
        this.balance = 1000;
        this.currentBet = 10;
        this.lastWin = 0;
        this.totalWins = 0;
        this.totalSpins = 0;
        this.gameHistory = [];
        this.saveGameState();
        this.updateDisplay();
    }
    
    // Random symbol generation with weighted probabilities
    getRandomSymbol() {
        // Weighted symbol selection for balanced gameplay
        const weights = {
            '🍀': 1,   // Rarest (200x)
            '💎': 2,   // Very rare (100x)
            '⭐': 3,   // Rare (75x)
            '🔔': 5,   // Uncommon (50x)
            '🍇': 8,   // Common (30x)
            '🍒': 12,  // Common (20x)
            '🍋': 15,  // Very common (15x)
            '🍊': 20   // Most common (10x)
        };
        
        const totalWeight = Object.values(weights).reduce((sum, weight) => sum + weight, 0);
        const random = Math.random() * totalWeight;
        
        let currentWeight = 0;
        for (const [symbol, weight] of Object.entries(weights)) {
            currentWeight += weight;
            if (random <= currentWeight) {
                return symbol;
            }
        }
        
        // Fallback to basic random
        return this.symbols[Math.floor(Math.random() * this.symbols.length)];
    }
    
    generateRandomResult() {
        return [
            this.getRandomSymbol(),
            this.getRandomSymbol(),
            this.getRandomSymbol()
        ];
    }
    
    // Enhanced win detection and analysis
    analyzeResult(result) {
        const [first, second, third] = result;
        
        // Check for exact three matches
        if (first === second && second === third) {
            return {
                type: 'three-of-a-kind',
                symbol: first,
                positions: [0, 1, 2],
                winLevel: this.getWinLevel(first + first + first)
            };
        }
        
        // Check for two-symbol matches
        if (first === second) {
            return {
                type: 'pair',
                symbol: first,
                positions: [0, 1],
                winLevel: this.getWinLevel(first + first)
            };
        }
        
        if (second === third) {
            return {
                type: 'pair',
                symbol: second,
                positions: [1, 2],
                winLevel: this.getWinLevel(second + second)
            };
        }
        
        if (first === third) {
            return {
                type: 'pair',
                symbol: first,
                positions: [0, 2],
                winLevel: this.getWinLevel(first + first)
            };
        }
        
        return {
            type: 'no-match',
            symbol: null,
            positions: [],
            winLevel: 'none'
        };
    }
    
    getWinLevel(combination) {
        const multiplier = this.payouts[combination] || 0;
        if (multiplier >= 100) return 'jackpot';
        if (multiplier >= 50) return 'mega';
        if (multiplier >= 20) return 'big';
        if (multiplier >= 10) return 'medium';
        if (multiplier > 0) return 'small';
        return 'none';
    }
    
    // UI updates
    updateDisplay() {
        // Update balance display
        const balanceElement = document.getElementById('balance');
        if (balanceElement) {
            balanceElement.textContent = this.balance.toFixed(0);
        }
        
        // Update last win display
        const lastWinElement = document.getElementById('last-win');
        if (lastWinElement) {
            lastWinElement.textContent = this.lastWin.toFixed(0);
        }
        
        // Update bet input
        const betInput = document.getElementById('bet-input');
        if (betInput) {
            betInput.value = this.currentBet;
            betInput.max = Math.min(100, this.balance);
        }
        
        // Update spin cost display
        const spinCostElement = document.getElementById('spin-cost');
        if (spinCostElement) {
            spinCostElement.textContent = this.currentBet;
        }
        
        // Update spin button state
        const spinButton = document.getElementById('spin-button');
        if (spinButton) {
            spinButton.disabled = !this.canAffordBet();
            if (!this.canAffordBet()) {
                spinButton.style.opacity = '0.5';
            } else {
                spinButton.style.opacity = '1';
            }
        }
    }
    
    // Debug/admin functions
    addDebugBalance(amount = 1000) {
        this.addToBalance(amount);
    }
    
    getDebugInfo() {
        return {
            state: this.getStats(),
            symbols: this.symbols,
            payouts: this.payouts,
            history: this.getHistory(5)
        };
    }
}

// Export for use in other modules
window.GameState = GameState;
