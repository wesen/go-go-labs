/**
 * UI Controller
 * Handles user interface interactions and updates
 */

class UIController {
    constructor(gameState, slotMachine, audioManager) {
        this.gameState = gameState;
        this.slotMachine = slotMachine;
        this.audioManager = audioManager;
        
        // Debouncing properties
        this.lastSpinTime = 0;
        this.spinCooldown = 500; // 500ms debounce
        this.pendingSpinTimeout = null;
        
        this.initializeUI();
        this.bindEvents();
    }
    
    initializeUI() {
        // Set initial UI state
        this.updateDisplay();
        this.updateBetPresets();
        this.initializeControls();
    }
    
    bindEvents() {
        // Spin button
        const spinButton = document.getElementById('spin-button');
        if (spinButton) {
            spinButton.addEventListener('click', () => {
                this.audioManager.playButtonClickSound();
                this.addButtonFeedback(spinButton);
                this.debouncedSpin();
            });
        }
        
        // Bet input
        const betInput = document.getElementById('bet-input');
        if (betInput) {
            betInput.addEventListener('input', (e) => {
                this.handleBetChange(e.target.value);
            });
            
            betInput.addEventListener('blur', () => {
                this.validateBetInput();
            });
        }
        
        // Bet preset buttons
        this.bindBetPresets();
        
        // Keyboard shortcuts
        this.bindKeyboardEvents();
        
        // Settings toggle (if implemented)
        this.bindSettingsEvents();
        
        // Balance click for debug (double-click to add money)
        const balanceElement = document.getElementById('balance');
        if (balanceElement) {
            let clickCount = 0;
            balanceElement.addEventListener('click', () => {
                clickCount++;
                setTimeout(() => {
                    if (clickCount === 2) {
                        this.handleDebugBalance();
                    }
                    clickCount = 0;
                }, 300);
            });
        }
    }
    
    bindBetPresets() {
        const betPresets = document.querySelectorAll('.bet-preset');
        betPresets.forEach(button => {
            button.addEventListener('click', (e) => {
                this.audioManager.playButtonClickSound();
                this.addButtonFeedback(e.target);
                const amount = parseInt(e.target.dataset.amount);
                this.setBetAmount(amount);
                this.updateBetPresets();
            });
        });
    }
    
    bindKeyboardEvents() {
        document.addEventListener('keydown', (e) => {
            // Prevent actions during spin
            if (this.slotMachine.isSpinning) return;
            
            switch (e.key) {
                case ' ':
                case 'Enter':
                    e.preventDefault();
                    if (this.gameState.canAffordBet()) {
                        this.debouncedSpin();
                    }
                    break;
                case '1':
                    this.setBetAmount(1);
                    break;
                case '2':
                    this.setBetAmount(5);
                    break;
                case '3':
                    this.setBetAmount(10);
                    break;
                case '4':
                    this.setBetAmount(25);
                    break;
                case '5':
                    this.setBetAmount(100);
                    break;
                case 'ArrowUp':
                    e.preventDefault();
                    this.increaseBet();
                    break;
                case 'ArrowDown':
                    e.preventDefault();
                    this.decreaseBet();
                    break;
                case 'm':
                case 'M':
                    this.toggleAudio();
                    break;
            }
        });
    }
    
    bindSettingsEvents() {
        // Audio toggle
        const audioToggle = document.getElementById('audio-toggle');
        if (audioToggle) {
            audioToggle.addEventListener('click', () => {
                this.toggleAudio();
            });
        }
        
        // Volume control
        const volumeSlider = document.getElementById('volume-slider');
        if (volumeSlider) {
            volumeSlider.addEventListener('input', (e) => {
                this.setVolume(e.target.value / 100);
            });
        }
    }
    
    // Debounced spin method to prevent rapid clicking
    debouncedSpin() {
        const now = Date.now();
        
        // Check if we're in cooldown period
        if (now - this.lastSpinTime < this.spinCooldown) {
            // Show visual feedback for too-fast clicking
            this.showCooldownFeedback();
            
            // Clear any pending timeout
            if (this.pendingSpinTimeout) {
                clearTimeout(this.pendingSpinTimeout);
            }
            
            // Schedule spin for after cooldown
            const remainingCooldown = this.spinCooldown - (now - this.lastSpinTime);
            this.pendingSpinTimeout = setTimeout(() => {
                this.pendingSpinTimeout = null;
                this.debouncedSpin();
            }, remainingCooldown);
            
            return;
        }
        
        // Clear any pending timeout
        if (this.pendingSpinTimeout) {
            clearTimeout(this.pendingSpinTimeout);
            this.pendingSpinTimeout = null;
        }
        
        // Update last spin time and execute spin
        this.lastSpinTime = now;
        this.handleSpin();
    }
    
    showCooldownFeedback() {
        const spinButton = document.getElementById('spin-button');
        if (spinButton) {
            spinButton.classList.add('cooldown-flash');
            setTimeout(() => {
                spinButton.classList.remove('cooldown-flash');
            }, 200);
        }
    }
    
    handleSpin() {
        if (this.slotMachine.isSpinning || !this.gameState.canAffordBet()) {
            return;
        }
        
        // Announce spin start for screen readers
        this.announceForScreenReader(`Spinning reels with $${this.gameState.getCurrentBet()} bet`);
        
        // Play spin sound
        this.audioManager.playSpinSound();
        
        // Update UI state
        this.updateSpinButtonState(false);
        
        // Start the spin
        this.slotMachine.spin().then(() => {
            const result = this.gameState.getLastWin();
            if (result > 0) {
                this.announceForScreenReader(`Spin complete! You won $${result}`);
            } else {
                this.announceForScreenReader('Spin complete. No win this time.');
            }
            this.updateDisplay();
            this.updateSpinButtonState(true);
        });
    }
    
    handleBetChange(value) {
        const betAmount = parseInt(value) || 1;
        this.setBetAmount(betAmount);
    }
    
    setBetAmount(amount) {
        const newBet = this.gameState.setCurrentBet(amount);
        this.updateDisplay();
        this.updateBetPresets();
        return newBet;
    }
    
    increaseBet() {
        const currentBet = this.gameState.getCurrentBet();
        const newBet = Math.min(currentBet + 1, Math.min(100, this.gameState.getBalance()));
        this.setBetAmount(newBet);
    }
    
    decreaseBet() {
        const currentBet = this.gameState.getCurrentBet();
        const newBet = Math.max(currentBet - 1, 1);
        this.setBetAmount(newBet);
    }
    
    validateBetInput() {
        const betInput = document.getElementById('bet-input');
        if (betInput) {
            const value = parseInt(betInput.value);
            if (isNaN(value) || value < 1) {
                betInput.value = 1;
                this.setBetAmount(1);
            } else if (value > this.gameState.getBalance()) {
                betInput.value = this.gameState.getBalance();
                this.setBetAmount(this.gameState.getBalance());
            }
        }
    }
    
    updateDisplay() {
        this.updateBalance();
        this.updateLastWin();
        this.updateBetDisplay();
        this.updateSpinButtonState(true);
    }
    
    updateBalance() {
        const balanceElement = document.getElementById('balance');
        if (balanceElement) {
            const balance = this.gameState.getBalance();
            balanceElement.textContent = balance.toLocaleString();
            
            // Add visual feedback for balance changes
            balanceElement.classList.add('bounce');
            setTimeout(() => {
                balanceElement.classList.remove('bounce');
            }, 600);
        }
    }
    
    updateLastWin() {
        const lastWinElement = document.getElementById('last-win');
        if (lastWinElement) {
            const lastWin = this.gameState.getLastWin();
            lastWinElement.textContent = lastWin.toLocaleString();
            
            if (lastWin > 0) {
                lastWinElement.parentElement.classList.add('pulse');
                setTimeout(() => {
                    lastWinElement.parentElement.classList.remove('pulse');
                }, 1000);
            }
        }
    }
    
    updateBetDisplay() {
        const betInput = document.getElementById('bet-input');
        const spinCost = document.getElementById('spin-cost');
        
        const currentBet = this.gameState.getCurrentBet();
        
        if (betInput && parseInt(betInput.value) !== currentBet) {
            betInput.value = currentBet;
        }
        
        if (spinCost) {
            spinCost.textContent = currentBet.toLocaleString();
        }
        
        // Update max bet limit
        if (betInput) {
            betInput.max = Math.min(100, this.gameState.getBalance());
        }
    }
    
    updateBetPresets() {
        const betPresets = document.querySelectorAll('.bet-preset');
        const currentBet = this.gameState.getCurrentBet();
        const balance = this.gameState.getBalance();
        
        betPresets.forEach(button => {
            const amount = parseInt(button.dataset.amount);
            
            // Update active state
            if (amount === currentBet) {
                button.classList.add('active');
            } else {
                button.classList.remove('active');
            }
            
            // Disable if can't afford
            if (amount > balance) {
                button.disabled = true;
                button.style.opacity = '0.5';
            } else {
                button.disabled = false;
                button.style.opacity = '1';
            }
        });
    }
    
    updateSpinButtonState(enabled) {
        const spinButton = document.getElementById('spin-button');
        if (!spinButton) return;
        
        const canSpin = enabled && this.gameState.canAffordBet() && !this.slotMachine.isSpinning;
        
        spinButton.disabled = !canSpin;
        
        if (canSpin) {
            spinButton.classList.remove('pulse');
            spinButton.style.opacity = '1';
        } else {
            spinButton.style.opacity = '0.6';
            if (this.slotMachine.isSpinning) {
                spinButton.classList.add('pulse');
            }
        }
        
        // Update button text
        const buttonText = spinButton.querySelector('.button-text');
        if (buttonText) {
            if (this.slotMachine.isSpinning) {
                buttonText.textContent = 'SPINNING...';
            } else if (!this.gameState.canAffordBet()) {
                buttonText.textContent = 'NO FUNDS';
            } else {
                buttonText.textContent = 'SPIN';
            }
        }
    }
    
    showMessage(message, type = 'normal', duration = 3000) {
        const messageDisplay = document.getElementById('message-display');
        if (!messageDisplay) return;
        
        messageDisplay.textContent = message;
        messageDisplay.className = `message-${type}`;
        
        // Add fade in effect
        messageDisplay.classList.add('fade-in');
        
        // Clear after duration
        setTimeout(() => {
            if (messageDisplay.textContent === message) {
                messageDisplay.classList.add('fade-out');
                setTimeout(() => {
                    messageDisplay.textContent = '';
                    messageDisplay.className = '';
                }, 500);
            }
        }, duration);
    }
    
    showWinMessage(amount, betAmount) {
        const multiplier = amount / betAmount;
        let message = '';
        let type = 'win';
        
        if (multiplier >= 100) {
            message = `🎰 MEGA JACKPOT! $${amount.toLocaleString()}! 🎰`;
            type = 'jackpot';
        } else if (multiplier >= 50) {
            message = `🎉 SUPER WIN! $${amount.toLocaleString()}! 🎉`;
            type = 'bigwin';
        } else if (multiplier >= 10) {
            message = `💰 BIG WIN! $${amount.toLocaleString()}! 💰`;
            type = 'bigwin';
        } else {
            message = `You won $${amount.toLocaleString()}!`;
            type = 'win';
        }
        
        this.showMessage(message, type, 4000);
    }
    
    // Audio controls
    toggleAudio() {
        const wasEnabled = this.audioManager.toggle();
        this.showMessage(
            wasEnabled ? 'Audio enabled' : 'Audio disabled',
            'info',
            1500
        );
        
        // Update audio button if it exists
        const audioButton = document.getElementById('audio-toggle');
        if (audioButton) {
            audioButton.textContent = wasEnabled ? '🔊' : '🔇';
        }
    }
    
    setVolume(volume) {
        this.audioManager.setVolume(volume);
    }
    
    // Debug functions
    handleDebugBalance() {
        if (process.env.NODE_ENV === 'development' || window.location.hostname === 'localhost') {
            this.gameState.addDebugBalance(1000);
            this.showMessage('Debug: Added $1000', 'info', 2000);
            this.updateDisplay();
        }
    }
    
    initializeControls() {
        // Set up any additional controls or widgets
        this.createKeyboardHelp();
        this.updateAudioControls();
    }
    
    createKeyboardHelp() {
        // Create a help tooltip for keyboard shortcuts
        const shortcuts = [
            'Space/Enter: Spin',
            '1-5: Quick bet amounts',
            '↑/↓: Adjust bet',
            'M: Toggle audio'
        ];
        
        // This could be implemented as a help button or tooltip
        console.log('Keyboard shortcuts:', shortcuts);
    }
    
    updateAudioControls() {
        const audioButton = document.getElementById('audio-toggle');
        if (audioButton) {
            audioButton.textContent = this.audioManager.isEnabled() ? '🔊' : '🔇';
        }
        
        const volumeSlider = document.getElementById('volume-slider');
        if (volumeSlider) {
            volumeSlider.value = this.audioManager.getVolume() * 100;
        }
    }
    
    // Enhanced visual feedback
    addButtonFeedback(button) {
        button.classList.add('button-press');
        setTimeout(() => {
            button.classList.remove('button-press');
        }, 150);
    }
    
    // Balance animation with more dramatic effects
    animateBalanceChange(oldBalance, newBalance) {
        const balanceElement = document.getElementById('balance');
        if (!balanceElement) return;
        
        const difference = newBalance - oldBalance;
        if (difference === 0) return;
        
        // Create floating text for change
        const changeText = document.createElement('div');
        changeText.className = 'balance-change';
        changeText.textContent = (difference > 0 ? '+' : '') + difference.toLocaleString();
        changeText.style.color = difference > 0 ? '#00ff00' : '#ff4444';
        
        const balanceContainer = balanceElement.parentElement;
        balanceContainer.style.position = 'relative';
        changeText.style.position = 'absolute';
        changeText.style.top = '-30px';
        changeText.style.left = '50%';
        changeText.style.transform = 'translateX(-50%)';
        changeText.style.fontSize = '1rem';
        changeText.style.fontWeight = 'bold';
        changeText.style.pointerEvents = 'none';
        changeText.style.animation = 'floatUp 2s ease-out forwards';
        
        balanceContainer.appendChild(changeText);
        
        setTimeout(() => {
            if (changeText.parentNode) {
                changeText.parentNode.removeChild(changeText);
            }
        }, 2000);
    }

    // Screen reader announcements
    announceForScreenReader(message) {
        const announceElement = document.getElementById('sr-announcements');
        if (announceElement) {
            announceElement.textContent = message;
            
            // Clear after 3 seconds to prepare for next announcement
            setTimeout(() => {
                announceElement.textContent = '';
            }, 3000);
        }
    }

    // Responsive UI adjustments
    handleResize() {
        // Handle responsive layout changes
        const gameContainer = document.querySelector('.game-container');
        if (gameContainer) {
            // Add mobile-specific classes if needed
            if (window.innerWidth < 768) {
                gameContainer.classList.add('mobile-layout');
            } else {
                gameContainer.classList.remove('mobile-layout');
            }
        }
    }
}

// Export for use in other modules
window.UIController = UIController;
