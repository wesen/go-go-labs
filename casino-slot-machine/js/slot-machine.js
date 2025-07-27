/**
 * Slot Machine Animation and Logic
 * Handles reel spinning, stopping, and visual effects
 */

class SlotMachine {
    constructor(gameState) {
        this.gameState = gameState;
        this.reels = [];
        this.isSpinning = false;
        this.spinSpeed = 100; // milliseconds per symbol change
        this.stopDelay = 500; // delay between reel stops
        this.symbols = gameState.symbols;
        
        // Accessibility mapping for emoji symbols
        this.symbolAccessibility = {
            '🍒': 'Cherry',
            '🍋': 'Lemon', 
            '🍊': 'Orange',
            '🔔': 'Bell',
            '🍇': 'Grapes',
            '⭐': 'Star',
            '💎': 'Diamond',
            '🍀': 'Lucky Clover'
        };
        
        this.initializeReels();
        this.bindEvents();
    }
    
    initializeReels() {
        // Initialize each reel
        for (let i = 1; i <= 3; i++) {
            const reelElement = document.getElementById(`reel-${i}`);
            if (reelElement) {
                const reel = {
                    element: reelElement,
                    symbols: [],
                    currentSymbol: 0,
                    spinning: false,
                    stopRequested: false,
                    finalSymbol: null
                };
                
                // Set initial symbols
                this.setInitialSymbols(reel);
                this.reels.push(reel);
            }
        }
    }
    
    setInitialSymbols(reel) {
        // Clear existing symbols
        reel.element.innerHTML = '';
        
        // Add initial symbols
        const initialSymbols = [
            this.gameState.getRandomSymbol(),
            this.gameState.getRandomSymbol(),
            this.gameState.getRandomSymbol()
        ];
        
        initialSymbols.forEach(symbol => {
            const symbolElement = this.createAccessibleSymbolElement(symbol);
            reel.element.appendChild(symbolElement);
        });
        
        reel.symbols = initialSymbols;
        reel.currentSymbol = 1; // Middle symbol is visible
    }
    
    // Create accessible symbol element with proper ARIA attributes
    createAccessibleSymbolElement(symbol) {
        const symbolElement = document.createElement('div');
        symbolElement.className = 'symbol';
        symbolElement.textContent = symbol;
        
        // Add accessibility attributes
        const altText = this.symbolAccessibility[symbol] || 'Slot symbol';
        symbolElement.setAttribute('aria-label', altText);
        symbolElement.setAttribute('role', 'img');
        symbolElement.setAttribute('title', altText);
        
        return symbolElement;
    }
    
    // Update symbol with accessibility
    updateSymbolWithAccessibility(element, symbol) {
        element.textContent = symbol;
        const altText = this.symbolAccessibility[symbol] || 'Slot symbol';
        element.setAttribute('aria-label', altText);
        element.setAttribute('title', altText);
    }
    
    bindEvents() {
        const spinButton = document.getElementById('spin-button');
        if (spinButton) {
            spinButton.addEventListener('click', () => {
                this.spin();
            });
        }
        
        // Bet input handling
        const betInput = document.getElementById('bet-input');
        if (betInput) {
            betInput.addEventListener('change', (e) => {
                const newBet = parseInt(e.target.value) || 1;
                this.gameState.setCurrentBet(newBet);
            });
        }
        
        // Bet preset buttons
        const betPresets = document.querySelectorAll('.bet-preset');
        betPresets.forEach(button => {
            button.addEventListener('click', (e) => {
                const amount = parseInt(e.target.dataset.amount);
                this.gameState.setCurrentBet(amount);
                
                // Update active preset
                betPresets.forEach(btn => btn.classList.remove('active'));
                e.target.classList.add('active');
            });
        });
    }
    
    async spin() {
        if (this.isSpinning || !this.gameState.canAffordBet()) {
            return;
        }
        
        this.isSpinning = true;
        this.updateSpinButton(false);
        this.clearWinEffects();
        
        try {
            // Generate final results
            const finalResults = this.gameState.generateRandomResult();
            
            // Start spinning all reels
            this.startSpinning();
            
            // Play enhanced audio feedback
            if (window.audioManager) {
                window.audioManager.playReelStartSound();
                window.audioManager.playElectricHumSound();
            }
            
            // Stop reels one by one with delays
            await this.stopReelsSequentially(finalResults);
            
            // Process the spin result
            const spinResult = this.gameState.processSpin(finalResults);
            const analysis = this.gameState.analyzeResult(finalResults);
            
            // Handle win/lose effects
            if (spinResult.isWin) {
                // Play enhanced win audio
                if (window.audioManager) {
                    window.audioManager.playWinSound(spinResult.payout, this.gameState.getCurrentBet(), analysis.winLevel);
                }
                await this.showWinEffects(spinResult.payout, analysis);
                // Add sparkle effect for wins
                if (window.audioManager && analysis && analysis.winLevel) {
                    window.audioManager.createSparkleEffect(
                        analysis.winLevel === 'jackpot' ? 5000 : 
                        analysis.winLevel === 'mega' ? 4000 : 3000
                    );
                }
            }
            
            this.showMessage(this.getSpinMessage(spinResult, analysis));
            
        } catch (error) {
            console.error('Spin error:', error);
            this.showMessage(error.message || 'Spin failed');
        } finally {
            this.isSpinning = false;
            this.updateSpinButton(true);
        }
    }
    
    startSpinning() {
        this.reels.forEach(reel => {
            reel.spinning = true;
            reel.stopRequested = false;
            reel.element.classList.add('spinning');
            this.animateReel(reel);
        });
    }
    
    animateReel(reel) {
        if (!reel.spinning) return;
        
        // Cycle through symbols
        const symbolElements = reel.element.querySelectorAll('.symbol');
        symbolElements.forEach((element, index) => {
            const newSymbol = this.gameState.getRandomSymbol();
            this.updateSymbolWithAccessibility(element, newSymbol);
        });
        
        // Continue animation
        setTimeout(() => {
            if (reel.spinning && !reel.stopRequested) {
                this.animateReel(reel);
            }
        }, this.spinSpeed);
    }
    
    async stopReelsSequentially(finalResults) {
        for (let i = 0; i < this.reels.length; i++) {
            await this.stopReel(this.reels[i], finalResults[i], i);
            await this.delay(this.stopDelay);
        }
    }
    
    async stopReel(reel, finalSymbol, reelIndex) {
        return new Promise(resolve => {
            reel.stopRequested = true;
            reel.finalSymbol = finalSymbol;
            
            // Add stopping class for animation
            reel.element.classList.add('stopping');
            reel.element.classList.remove('spinning');
            
            // Play reel stop sound with pitch variation
            if (window.audioManager) {
                window.audioManager.playReelStopSoundWithPitch(reelIndex);
            }
            
            setTimeout(() => {
                reel.spinning = false;
                
                // Set final symbol
                const symbolElements = reel.element.querySelectorAll('.symbol');
                if (symbolElements[1]) { // Middle symbol
                    this.updateSymbolWithAccessibility(symbolElements[1], finalSymbol);
                }
                
                reel.element.classList.remove('stopping');
                
                // Add bounce effect
                reel.element.classList.add('bounce');
                setTimeout(() => {
                    reel.element.classList.remove('bounce');
                    resolve();
                }, 600);
                
            }, 500);
        });
    }
    
    async showWinEffects(winAmount, analysis) {
        const reelsContainer = document.querySelector('.reels-container');
        const payline = document.querySelector('.payline');
        const slotMachine = document.querySelector('.slot-machine');
        
        // Add win classes
        if (payline) payline.classList.add('winning');
        if (slotMachine) slotMachine.classList.add('celebration');
        
        // Add winning symbol effects to specific positions
        if (analysis && analysis.positions) {
            analysis.positions.forEach(pos => {
                const symbolElement = document.querySelector(`.reel:nth-child(${pos + 1}) .symbol:nth-child(2)`);
                if (symbolElement) {
                    symbolElement.classList.add('winning');
                    symbolElement.classList.add(`win-${analysis.winLevel}`);
                }
            });
        }
        
        // Special effects based on win level
        if (analysis && analysis.winLevel) {
            if (reelsContainer) reelsContainer.classList.add(`effect-${analysis.winLevel}`);
            
            switch (analysis.winLevel) {
                case 'jackpot':
                    this.showMessage(`🎰 JACKPOT! 🎰`, 'jackpot');
                    this.createLightningEffect();
                    this.createCoinCascade(50);
                    break;
                case 'mega':
                    this.showMessage(`💰 MEGA WIN! 💰`, 'mega');
                    this.createCoinCascade(30);
                    this.createParticleEffect('star', 25);
                    break;
                case 'big':
                    this.showMessage(`🎉 BIG WIN! 🎉`, 'big');
                    this.createCoinCascade(20);
                    this.createParticleEffect('coin', 15);
                    break;
                case 'medium':
                    this.createParticleEffect('star', 10);
                    break;
                case 'small':
                    this.createParticleEffect('coin', 5);
                    break;
            }
        }
        
        // Flash effect duration varies by win level
        const duration = analysis && analysis.winLevel === 'jackpot' ? 4000 : 
                        analysis && analysis.winLevel === 'mega' ? 3000 : 2000;
        await this.delay(duration);
    }
    
    clearWinEffects() {
        // Remove all win-related classes
        const elements = [
            '.payline',
            '.slot-machine',
            '.reels-container',
            '.symbol'
        ];
        
        elements.forEach(selector => {
            const elements = document.querySelectorAll(selector);
            elements.forEach(element => {
                element.classList.remove('winning', 'celebration', 'jackpot');
            });
        });
    }
    
    updateSpinButton(enabled) {
        const spinButton = document.getElementById('spin-button');
        if (spinButton) {
            spinButton.disabled = !enabled || !this.gameState.canAffordBet();
            
            if (!enabled) {
                spinButton.classList.add('pulse');
            } else {
                spinButton.classList.remove('pulse');
            }
        }
    }
    
    showMessage(message, type = 'normal') {
        const messageDisplay = document.getElementById('message-display');
        if (messageDisplay) {
            messageDisplay.textContent = message;
            messageDisplay.className = `message-${type}`;
            
            // Add fade in effect
            messageDisplay.classList.add('fade-in');
            
            // Clear message after delay
            setTimeout(() => {
                if (messageDisplay.textContent === message) {
                    messageDisplay.textContent = '';
                    messageDisplay.className = '';
                }
            }, 3000);
        }
    }
    
    getSpinMessage(spinResult, analysis) {
        if (spinResult.payout === 0) {
            const messages = [
                'Try again!',
                'Better luck next time!',
                'Keep spinning!',
                'Fortune favors the bold!',
                'Next spin could be the one!'
            ];
            return messages[Math.floor(Math.random() * messages.length)];
        }
        
        // Use analysis for more specific messages
        if (analysis && analysis.winLevel) {
            switch (analysis.winLevel) {
                case 'jackpot':
                    return `🍀 INCREDIBLE! ${analysis.symbol}${analysis.symbol}${analysis.symbol} = $${spinResult.payout}! 🍀`;
                case 'mega':
                    return `💎 FANTASTIC! ${analysis.symbol}${analysis.symbol}${analysis.symbol} = $${spinResult.payout}! 💎`;
                case 'big':
                    return `⭐ AWESOME! ${analysis.symbol}${analysis.symbol}${analysis.symbol} = $${spinResult.payout}! ⭐`;
                case 'medium':
                    return `🔔 NICE! ${analysis.symbol}${analysis.symbol}${analysis.symbol} = $${spinResult.payout}! 🔔`;
                case 'small':
                    return `${analysis.symbol} You won $${spinResult.payout}! ${analysis.symbol}`;
            }
        }
        
        // Fallback to generic messages
        if (spinResult.payout >= this.gameState.getCurrentBet() * 50) {
            return `🎰 MASSIVE WIN! $${spinResult.payout}! 🎰`;
        } else if (spinResult.payout >= this.gameState.getCurrentBet() * 10) {
            return `🎉 BIG WIN! $${spinResult.payout}! 🎉`;
        } else {
            return `You won $${spinResult.payout}!`;
        }
    }
    
    // Visual effects methods
    createParticleEffect(type = 'coin', count = 10) {
        const container = document.querySelector('.slot-machine-frame');
        if (!container) return;
        
        // Create particles container if it doesn't exist
        let particlesContainer = container.querySelector('.particles-container');
        if (!particlesContainer) {
            particlesContainer = document.createElement('div');
            particlesContainer.className = 'particles-container';
            container.appendChild(particlesContainer);
        }
        
        for (let i = 0; i < count; i++) {
            setTimeout(() => {
                const particle = document.createElement('div');
                particle.className = `particle ${type}`;
                
                // Random position across the top
                const leftPos = Math.random() * 80 + 10; // 10% to 90%
                particle.style.left = `${leftPos}%`;
                particle.style.top = '0';
                
                // Random drift
                const drift = (Math.random() - 0.5) * 100;
                particle.style.setProperty('--drift', `${drift}px`);
                
                // Set content based on type
                if (type === 'star') {
                    particle.textContent = ['⭐', '✨', '💫', '🌟'][Math.floor(Math.random() * 4)];
                } else if (type === 'coin') {
                    particle.textContent = '$';
                }
                
                particlesContainer.appendChild(particle);
                
                // Remove particle after animation
                setTimeout(() => {
                    if (particle.parentNode) {
                        particle.parentNode.removeChild(particle);
                    }
                }, 2000);
            }, i * 100);
        }
    }
    
    createCoinCascade(count = 20) {
        const container = document.querySelector('.slot-machine-frame');
        if (!container) return;
        
        // Create cascade container if it doesn't exist
        let cascadeContainer = container.querySelector('.coin-cascade');
        if (!cascadeContainer) {
            cascadeContainer = document.createElement('div');
            cascadeContainer.className = 'coin-cascade';
            container.appendChild(cascadeContainer);
        }
        
        for (let i = 0; i < count; i++) {
            setTimeout(() => {
                const coin = document.createElement('div');
                coin.className = 'falling-coin';
                coin.textContent = '$';
                
                // Random position across the top
                const leftPos = Math.random() * 90 + 5; // 5% to 95%
                coin.style.left = `${leftPos}%`;
                coin.style.top = '-30px';
                
                // Slight random delay for more natural effect
                coin.style.animationDelay = `${Math.random() * 0.3}s`;
                
                cascadeContainer.appendChild(coin);
                
                // Remove coin after animation
                setTimeout(() => {
                    if (coin.parentNode) {
                        coin.parentNode.removeChild(coin);
                    }
                }, 2000);
            }, i * 50);
        }
    }
    
    createLightningEffect() {
        const container = document.querySelector('.slot-machine-frame');
        if (!container) return;
        
        const lightning = document.createElement('div');
        lightning.className = 'lightning-effect';
        container.appendChild(lightning);
        
        // Remove after animation
        setTimeout(() => {
            if (lightning.parentNode) {
                lightning.parentNode.removeChild(lightning);
            }
        }, 600);
    }
    
    // Utility methods
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    getCurrentResults() {
        return this.reels.map(reel => {
            const symbolElement = reel.element.querySelector('.symbol:nth-child(2)');
            return symbolElement ? symbolElement.textContent : '?';
        });
    }
    
    // Debug methods
    forceResult(symbols) {
        if (this.isSpinning) return false;
        
        symbols.forEach((symbol, index) => {
            if (this.reels[index]) {
                const symbolElement = this.reels[index].element.querySelector('.symbol:nth-child(2)');
                if (symbolElement) {
                    this.updateSymbolWithAccessibility(symbolElement, symbol);
                }
            }
        });
        
        return true;
    }
    
    getDebugInfo() {
        return {
            isSpinning: this.isSpinning,
            currentResults: this.getCurrentResults(),
            reelStates: this.reels.map(reel => ({
                spinning: reel.spinning,
                stopRequested: reel.stopRequested,
                symbols: reel.symbols
            }))
        };
    }
}

// Export for use in other modules
window.SlotMachine = SlotMachine;
