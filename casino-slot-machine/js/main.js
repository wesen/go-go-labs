/**
 * Main Application
 * Initializes and coordinates all game components
 */

class CasinoSlotMachine {
    constructor() {
        this.gameState = null;
        this.slotMachine = null;
        this.audioManager = null;
        this.uiController = null;
        this.initialized = false;
        
        this.initialize();
    }
    
    async initialize() {
        try {
            console.log('🎰 Initializing Casino Slot Machine...');
            
            // Wait for DOM to be ready
            if (document.readyState === 'loading') {
                await new Promise(resolve => {
                    document.addEventListener('DOMContentLoaded', resolve);
                });
            }
            
            // Initialize core components
            this.initializeComponents();
            
            // Set up error handling
            this.setupErrorHandling();
            
            // Set up performance monitoring
            this.setupPerformanceMonitoring();
            
            // Set up resize handling
            this.setupResizeHandling();
            
            // Mark as initialized
            this.initialized = true;
            
            console.log('✅ Casino Slot Machine initialized successfully');
            
            // Show welcome message
            this.showWelcomeMessage();
            
        } catch (error) {
            console.error('❌ Failed to initialize Casino Slot Machine:', error);
            this.showErrorMessage('Failed to initialize game. Please refresh the page.');
        }
    }
    
    initializeComponents() {
        // Initialize game state (must be first)
        this.gameState = new GameState();
        console.log('✅ Game State initialized');
        
        // Initialize audio manager
        this.audioManager = new AudioManager();
        console.log('✅ Audio Manager initialized');
        
        // Initialize slot machine
        this.slotMachine = new SlotMachine(this.gameState);
        console.log('✅ Slot Machine initialized');
        
        // Initialize UI controller (must be last)
        this.uiController = new UIController(
            this.gameState,
            this.slotMachine,
            this.audioManager
        );
        console.log('✅ UI Controller initialized');
        
        // Set up cross-component communication
        this.setupComponentCommunication();
    }
    
    setupComponentCommunication() {
        // Override slot machine spin method to include audio
        const originalSpin = this.slotMachine.spin.bind(this.slotMachine);
        this.slotMachine.spin = async () => {
            this.audioManager.playSpinSound();
            const result = await originalSpin();
            
            // Play appropriate win sound if there was a win
            if (this.gameState.getLastWin() > 0) {
                this.audioManager.playWinSound(
                    this.gameState.getLastWin(),
                    this.gameState.getCurrentBet()
                );
            }
            
            return result;
        };
        
        // Add reel stop sounds
        const originalStopReel = this.slotMachine.stopReel.bind(this.slotMachine);
        this.slotMachine.stopReel = async (reel, finalSymbol, reelIndex) => {
            const result = await originalStopReel(reel, finalSymbol, reelIndex);
            this.audioManager.playReelStopSound();
            return result;
        };
    }
    
    setupErrorHandling() {
        // Global error handler
        window.addEventListener('error', (event) => {
            console.error('Global error:', event.error);
            this.handleError(event.error);
        });
        
        // Promise rejection handler
        window.addEventListener('unhandledrejection', (event) => {
            console.error('Unhandled promise rejection:', event.reason);
            this.handleError(event.reason);
        });
    }
    
    setupPerformanceMonitoring() {
        // Monitor frame rate and performance
        let frameCount = 0;
        let lastTime = performance.now();
        
        const measurePerformance = () => {
            frameCount++;
            const currentTime = performance.now();
            
            if (currentTime - lastTime >= 1000) {
                const fps = Math.round(frameCount * 1000 / (currentTime - lastTime));
                
                // Log performance issues
                if (fps < 30) {
                    console.warn(`⚠️ Low FPS detected: ${fps}`);
                }
                
                frameCount = 0;
                lastTime = currentTime;
            }
            
            requestAnimationFrame(measurePerformance);
        };
        
        requestAnimationFrame(measurePerformance);
    }
    
    setupResizeHandling() {
        let resizeTimeout;
        
        window.addEventListener('resize', () => {
            clearTimeout(resizeTimeout);
            resizeTimeout = setTimeout(() => {
                this.handleResize();
            }, 250);
        });
        
        // Initial resize
        this.handleResize();
    }
    
    handleResize() {
        if (this.uiController) {
            this.uiController.handleResize();
        }
        
        // Adjust game layout for different screen sizes
        this.adjustLayoutForScreen();
    }
    
    adjustLayoutForScreen() {
        const gameContainer = document.querySelector('.game-container');
        if (!gameContainer) return;
        
        const width = window.innerWidth;
        const height = window.innerHeight;
        
        // Add size classes for CSS targeting
        gameContainer.classList.remove('small-screen', 'medium-screen', 'large-screen');
        
        if (width < 768) {
            gameContainer.classList.add('small-screen');
        } else if (width < 1200) {
            gameContainer.classList.add('medium-screen');
        } else {
            gameContainer.classList.add('large-screen');
        }
        
        // Handle landscape mode on mobile
        if (width < 768 && width > height) {
            gameContainer.classList.add('mobile-landscape');
        } else {
            gameContainer.classList.remove('mobile-landscape');
        }
    }
    
    handleError(error) {
        // Don't spam error messages
        if (this.lastErrorTime && Date.now() - this.lastErrorTime < 5000) {
            return;
        }
        this.lastErrorTime = Date.now();
        
        // Show user-friendly error message
        this.showErrorMessage('Something went wrong. The game will continue to work.');
        
        // Try to recover
        this.attemptRecovery();
    }
    
    attemptRecovery() {
        try {
            // Refresh UI state
            if (this.uiController) {
                this.uiController.updateDisplay();
            }
            
            // Clear any stuck animations
            this.clearAnimations();
            
            console.log('🔄 Attempted error recovery');
        } catch (recoveryError) {
            console.error('❌ Recovery failed:', recoveryError);
        }
    }
    
    clearAnimations() {
        // Remove any stuck animation classes
        const animatedElements = document.querySelectorAll('.spinning, .stopping, .winning, .celebration');
        animatedElements.forEach(element => {
            element.classList.remove('spinning', 'stopping', 'winning', 'celebration');
        });
    }
    
    showWelcomeMessage() {
        if (this.uiController) {
            const messages = [
                'Welcome to the Casino! 🎰',
                'Good luck and have fun! 🍀',
                'Try your luck! 💰',
                'Feeling lucky? 🎲'
            ];
            
            const randomMessage = messages[Math.floor(Math.random() * messages.length)];
            this.uiController.showMessage(randomMessage, 'info', 3000);
        }
    }
    
    showErrorMessage(message) {
        if (this.uiController) {
            this.uiController.showMessage(`⚠️ ${message}`, 'error', 5000);
        } else {
            // Fallback if UI controller isn't available
            alert(message);
        }
    }
    
    // Public API methods
    getGameState() {
        return this.gameState;
    }
    
    getSlotMachine() {
        return this.slotMachine;
    }
    
    getAudioManager() {
        return this.audioManager;
    }
    
    getUIController() {
        return this.uiController;
    }
    
    // Debug methods
    enableDebugMode() {
        window.casino = this;
        window.gameState = this.gameState;
        window.slotMachine = this.slotMachine;
        window.audioManager = this.audioManager;
        window.uiController = this.uiController;
        
        console.log('🐛 Debug mode enabled. Access components via window.casino');
        console.log('Available: window.gameState, window.slotMachine, window.audioManager, window.uiController');
    }
    
    getDebugInfo() {
        return {
            initialized: this.initialized,
            gameState: this.gameState?.getDebugInfo(),
            slotMachine: this.slotMachine?.getDebugInfo(),
            audio: {
                enabled: this.audioManager?.isEnabled(),
                volume: this.audioManager?.getVolume()
            },
            screen: {
                width: window.innerWidth,
                height: window.innerHeight,
                devicePixelRatio: window.devicePixelRatio
            },
            performance: {
                memory: performance.memory ? {
                    used: Math.round(performance.memory.usedJSHeapSize / 1024 / 1024),
                    total: Math.round(performance.memory.totalJSHeapSize / 1024 / 1024),
                    limit: Math.round(performance.memory.jsHeapSizeLimit / 1024 / 1024)
                } : 'Not available'
            }
        };
    }
    
    // Lifecycle methods
    destroy() {
        // Clean up event listeners
        window.removeEventListener('error', this.handleError);
        window.removeEventListener('unhandledrejection', this.handleError);
        
        // Save game state
        if (this.gameState) {
            this.gameState.saveGameState();
        }
        
        // Clean up audio context
        if (this.audioManager && this.audioManager.context) {
            this.audioManager.context.close();
        }
        
        console.log('🧹 Casino Slot Machine cleaned up');
    }
}

// Initialize the game when the script loads
document.addEventListener('DOMContentLoaded', () => {
    window.casinoGame = new CasinoSlotMachine();
    
    // Enable debug mode in development
    if (window.location.hostname === 'localhost' || 
        window.location.hostname === '127.0.0.1' ||
        window.location.search.includes('debug=true')) {
        window.casinoGame.enableDebugMode();
    }
});

// Handle page unload
window.addEventListener('beforeunload', () => {
    if (window.casinoGame) {
        window.casinoGame.destroy();
    }
});

// Export for manual initialization if needed
window.CasinoSlotMachine = CasinoSlotMachine;
