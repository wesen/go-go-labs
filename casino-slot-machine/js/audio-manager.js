/**
 * Audio Manager
 * Handles sound effects and background music for the slot machine
 */

class AudioManager {
    constructor() {
        this.sounds = {};
        this.enabled = true;
        this.volume = 0.7;
        this.musicVolume = 0.3;
        this.context = null;
        this.audioPermissionGranted = false;
        this.permissionPromptShown = false;
        
        this.initializeAudio();
        this.createSynthSounds();
        this.loadSettings();
        this.setupAudioPermissionPrompt();
    }
    
    initializeAudio() {
        // Check for Web Audio API support
        try {
            window.AudioContext = window.AudioContext || window.webkitAudioContext;
            this.context = new AudioContext();
        } catch (error) {
            console.warn('Web Audio API not supported, falling back to basic audio');
        }
    }
    
    createSynthSounds() {
        if (!this.context) return;
        
        // Create synthetic sound effects using Web Audio API
        this.sounds = {
            spin: () => this.createSpinSound(),
            win: () => this.createWinSound(),
            bigWin: () => this.createBigWinSound(),
            jackpot: () => this.createJackpotSound(),
            reelStop: () => this.createReelStopSound(),
            buttonClick: () => this.createButtonClickSound(),
            coinDrop: () => this.createCoinDropSound(),
            background: () => this.createBackgroundMusic(),
            reelStart: () => this.createReelStartSound(),
            anticipation: () => this.createAnticipationSound(),
            celebration: () => this.createCelebrationSound(),
            electricHum: () => this.createElectricHumSound(),
            coinCascade: () => this.createCoinCascadeSound()
        };
    }
    
    createSpinSound() {
        if (!this.context) return;
        
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const filter = this.context.createBiquadFilter();
        
        oscillator.connect(filter);
        filter.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        // Spinning reel sound - low frequency with modulation
        oscillator.type = 'sawtooth';
        oscillator.frequency.setValueAtTime(60, this.context.currentTime);
        oscillator.frequency.exponentialRampToValueAtTime(120, this.context.currentTime + 0.1);
        oscillator.frequency.exponentialRampToValueAtTime(80, this.context.currentTime + 0.2);
        
        filter.type = 'lowpass';
        filter.frequency.setValueAtTime(800, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.3, this.context.currentTime + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 2);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 2);
    }
    
    createWinSound() {
        if (!this.context) return;
        
        // Happy ascending chord progression
        const notes = [261.63, 329.63, 392.00, 523.25]; // C, E, G, C (major chord)
        
        notes.forEach((freq, index) => {
            const oscillator = this.context.createOscillator();
            const gainNode = this.context.createGain();
            
            oscillator.connect(gainNode);
            gainNode.connect(this.context.destination);
            
            oscillator.type = 'sine';
            oscillator.frequency.setValueAtTime(freq, this.context.currentTime);
            
            const startTime = this.context.currentTime + (index * 0.1);
            gainNode.gain.setValueAtTime(0, startTime);
            gainNode.gain.linearRampToValueAtTime(this.volume * 0.4, startTime + 0.01);
            gainNode.gain.exponentialRampToValueAtTime(0.01, startTime + 0.5);
            
            oscillator.start(startTime);
            oscillator.stop(startTime + 0.5);
        });
    }
    
    createBigWinSound() {
        if (!this.context) return;
        
        // More elaborate win sound with multiple layers
        for (let i = 0; i < 3; i++) {
            setTimeout(() => {
                this.createWinSound();
                
                // Add sparkle effect
                const sparkle = this.context.createOscillator();
                const sparkleGain = this.context.createGain();
                
                sparkle.connect(sparkleGain);
                sparkleGain.connect(this.context.destination);
                
                sparkle.type = 'sine';
                sparkle.frequency.setValueAtTime(1000 + (i * 200), this.context.currentTime);
                
                sparkleGain.gain.setValueAtTime(0, this.context.currentTime);
                sparkleGain.gain.linearRampToValueAtTime(this.volume * 0.2, this.context.currentTime + 0.01);
                sparkleGain.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.3);
                
                sparkle.start(this.context.currentTime);
                sparkle.stop(this.context.currentTime + 0.3);
            }, i * 150);
        }
    }
    
    createJackpotSound() {
        if (!this.context) return;
        
        // Epic jackpot fanfare
        const fanfare = [
            { freq: 523.25, time: 0.0 },    // C
            { freq: 659.25, time: 0.2 },    // E
            { freq: 783.99, time: 0.4 },    // G
            { freq: 1046.50, time: 0.6 },   // C (octave)
            { freq: 1318.51, time: 0.8 },   // E (octave)
        ];
        
        fanfare.forEach(note => {
            const oscillator = this.context.createOscillator();
            const gainNode = this.context.createGain();
            const reverb = this.context.createConvolver();
            
            oscillator.connect(gainNode);
            gainNode.connect(reverb);
            reverb.connect(this.context.destination);
            
            oscillator.type = 'triangle';
            oscillator.frequency.setValueAtTime(note.freq, this.context.currentTime + note.time);
            
            gainNode.gain.setValueAtTime(0, this.context.currentTime + note.time);
            gainNode.gain.linearRampToValueAtTime(this.volume * 0.6, this.context.currentTime + note.time + 0.01);
            gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + note.time + 1);
            
            oscillator.start(this.context.currentTime + note.time);
            oscillator.stop(this.context.currentTime + note.time + 1);
        });
    }
    
    createReelStopSound() {
        if (!this.context) return;
        
        // Mechanical stop sound
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const filter = this.context.createBiquadFilter();
        
        oscillator.connect(filter);
        filter.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        oscillator.type = 'square';
        oscillator.frequency.setValueAtTime(200, this.context.currentTime);
        oscillator.frequency.exponentialRampToValueAtTime(50, this.context.currentTime + 0.1);
        
        filter.type = 'highpass';
        filter.frequency.setValueAtTime(100, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.3, this.context.currentTime + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.15);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 0.15);
    }
    
    createButtonClickSound() {
        if (!this.context) return;
        
        // Quick click sound
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        oscillator.type = 'square';
        oscillator.frequency.setValueAtTime(800, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.2, this.context.currentTime + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.05);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 0.05);
    }
    
    createCoinDropSound() {
        if (!this.context) return;
        
        // Coin drop/collect sound
        const frequencies = [1000, 800, 600, 400];
        
        frequencies.forEach((freq, index) => {
            setTimeout(() => {
                const oscillator = this.context.createOscillator();
                const gainNode = this.context.createGain();
                
                oscillator.connect(gainNode);
                gainNode.connect(this.context.destination);
                
                oscillator.type = 'sine';
                oscillator.frequency.setValueAtTime(freq, this.context.currentTime);
                
                gainNode.gain.setValueAtTime(0, this.context.currentTime);
                gainNode.gain.linearRampToValueAtTime(this.volume * 0.3, this.context.currentTime + 0.01);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.1);
                
                oscillator.start(this.context.currentTime);
                oscillator.stop(this.context.currentTime + 0.1);
            }, index * 20);
        });
    }
    
    createReelStartSound() {
        if (!this.context) return;
        
        // Quick rising tone to indicate reel starting
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const filter = this.context.createBiquadFilter();
        
        oscillator.connect(filter);
        filter.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        oscillator.type = 'sine';
        oscillator.frequency.setValueAtTime(220, this.context.currentTime);
        oscillator.frequency.exponentialRampToValueAtTime(440, this.context.currentTime + 0.1);
        
        filter.type = 'lowpass';
        filter.frequency.setValueAtTime(1000, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.2, this.context.currentTime + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.1);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 0.1);
    }
    
    createAnticipationSound() {
        if (!this.context) return;
        
        // Building tension sound
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const lfo = this.context.createOscillator();
        const lfoGain = this.context.createGain();
        
        lfo.connect(lfoGain);
        lfoGain.connect(oscillator.frequency);
        oscillator.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        oscillator.type = 'triangle';
        oscillator.frequency.setValueAtTime(150, this.context.currentTime);
        
        lfo.type = 'sine';
        lfo.frequency.setValueAtTime(8, this.context.currentTime);
        lfoGain.gain.setValueAtTime(20, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.1, this.context.currentTime + 0.1);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.3, this.context.currentTime + 1);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 1.5);
        
        lfo.start(this.context.currentTime);
        oscillator.start(this.context.currentTime);
        lfo.stop(this.context.currentTime + 1.5);
        oscillator.stop(this.context.currentTime + 1.5);
    }
    
    createCelebrationSound() {
        if (!this.context) return;
        
        // Rapid ascending arpeggios
        const notes = [523.25, 659.25, 783.99, 1046.50, 1318.51]; // C major pentatonic
        
        for (let i = 0; i < 3; i++) {
            notes.forEach((freq, index) => {
                setTimeout(() => {
                    const oscillator = this.context.createOscillator();
                    const gainNode = this.context.createGain();
                    
                    oscillator.connect(gainNode);
                    gainNode.connect(this.context.destination);
                    
                    oscillator.type = 'square';
                    oscillator.frequency.setValueAtTime(freq, this.context.currentTime);
                    
                    gainNode.gain.setValueAtTime(0, this.context.currentTime);
                    gainNode.gain.linearRampToValueAtTime(this.volume * 0.3, this.context.currentTime + 0.01);
                    gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.15);
                    
                    oscillator.start(this.context.currentTime);
                    oscillator.stop(this.context.currentTime + 0.15);
                }, (i * 300) + (index * 50));
            });
        }
    }
    
    createElectricHumSound() {
        if (!this.context) return;
        
        // Electric/mechanical hum for machine ambience
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const filter = this.context.createBiquadFilter();
        
        oscillator.connect(filter);
        filter.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        oscillator.type = 'sawtooth';
        oscillator.frequency.setValueAtTime(60, this.context.currentTime);
        
        filter.type = 'lowpass';
        filter.frequency.setValueAtTime(200, this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.05, this.context.currentTime + 0.1);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.05, this.context.currentTime + 3);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 3.5);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 3.5);
    }
    
    createCoinCascadeSound() {
        if (!this.context) return;
        
        // Multiple overlapping coin sounds for cascade effect
        for (let i = 0; i < 20; i++) {
            setTimeout(() => {
                const frequencies = [800, 1000, 1200];
                const freq = frequencies[Math.floor(Math.random() * frequencies.length)];
                
                const oscillator = this.context.createOscillator();
                const gainNode = this.context.createGain();
                
                oscillator.connect(gainNode);
                gainNode.connect(this.context.destination);
                
                oscillator.type = 'sine';
                oscillator.frequency.setValueAtTime(freq, this.context.currentTime);
                oscillator.frequency.exponentialRampToValueAtTime(freq * 0.5, this.context.currentTime + 0.1);
                
                gainNode.gain.setValueAtTime(0, this.context.currentTime);
                gainNode.gain.linearRampToValueAtTime(this.volume * 0.2, this.context.currentTime + 0.01);
                gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.15);
                
                oscillator.start(this.context.currentTime);
                oscillator.stop(this.context.currentTime + 0.15);
            }, i * 50);
        }
    }
    
    // Main sound playing methods
    playSound(soundName) {
        if (!this.enabled || !this.sounds[soundName]) return;
        
        try {
            // Check if we need to request audio permission
            if (!this.audioPermissionGranted && !this.permissionPromptShown) {
                this.showAudioPermissionPrompt();
                return;
            }
            
            // Resume audio context if suspended (required by some browsers)
            if (this.context && this.context.state === 'suspended') {
                this.context.resume().then(() => {
                    this.audioPermissionGranted = true;
                    this.sounds[soundName]();
                });
            } else {
                this.sounds[soundName]();
            }
        } catch (error) {
            console.warn('Error playing sound:', soundName, error);
        }
    }
    
    playSpinSound() {
        this.playSound('spin');
    }
    
    playWinSound(amount, betAmount, winLevel = null) {
        const multiplier = amount / betAmount;
        
        if (winLevel === 'jackpot' || multiplier >= 100) {
            this.playSound('jackpot');
            setTimeout(() => this.playSound('celebration'), 500);
            setTimeout(() => this.playSound('coinCascade'), 800);
            // Add screen flash effect
            this.createScreenFlash('jackpot');
        } else if (winLevel === 'mega' || winLevel === 'big' || multiplier >= 20) {
            this.playSound('bigWin');
            setTimeout(() => this.playSound('celebration'), 300);
            setTimeout(() => this.playSound('coinDrop'), 600);
            // Add screen flash effect
            this.createScreenFlash('mega');
        } else {
            this.playSound('win');
            setTimeout(() => this.playSound('coinDrop'), 400);
        }
    }
    
    playReelStopSound() {
        this.playSound('reelStop');
    }
    
    playButtonClickSound() {
        this.playSound('buttonClick');
    }
    
    playReelStartSound() {
        this.playSound('reelStart');
    }
    
    playAnticipationSound() {
        this.playSound('anticipation');
    }
    
    playCelebrationSound() {
        this.playSound('celebration');
    }
    
    playElectricHumSound() {
        this.playSound('electricHum');
    }
    
    playCoinCascadeSound() {
        this.playSound('coinCascade');
    }
    
    // Settings management
    setEnabled(enabled) {
        this.enabled = enabled;
        this.saveSettings();
    }
    
    setVolume(volume) {
        this.volume = Math.max(0, Math.min(1, volume));
        this.saveSettings();
    }
    
    isEnabled() {
        return this.enabled;
    }
    
    getVolume() {
        return this.volume;
    }
    
    toggle() {
        this.setEnabled(!this.enabled);
        return this.enabled;
    }
    
    // Persistence
    saveSettings() {
        try {
            const settings = {
                enabled: this.enabled,
                volume: this.volume,
                musicVolume: this.musicVolume
            };
            localStorage.setItem('slotMachineAudioSettings', JSON.stringify(settings));
        } catch (error) {
            console.warn('Could not save audio settings:', error);
        }
    }
    
    loadSettings() {
        try {
            const saved = localStorage.getItem('slotMachineAudioSettings');
            if (saved) {
                const settings = JSON.parse(saved);
                this.enabled = settings.enabled !== false; // Default to true
                this.volume = settings.volume || 0.7;
                this.musicVolume = settings.musicVolume || 0.3;
            }
        } catch (error) {
            console.warn('Could not load audio settings:', error);
        }
    }
    
    // Visual effects integration
    createScreenFlash(type) {
        const flash = document.createElement('div');
        flash.className = `screen-flash ${type}`;
        document.body.appendChild(flash);
        
        setTimeout(() => {
            if (flash.parentNode) {
                flash.parentNode.removeChild(flash);
            }
        }, type === 'jackpot' ? 500 : 450);
    }
    
    // Add ambient sparkles during wins
    createSparkleEffect(duration = 3000) {
        const sparkleContainer = document.createElement('div');
        sparkleContainer.className = 'sparkle-bg';
        
        const slotMachine = document.querySelector('.slot-machine-frame');
        if (slotMachine) {
            slotMachine.appendChild(sparkleContainer);
            
            // Create multiple sparkles
            for (let i = 0; i < 15; i++) {
                setTimeout(() => {
                    const sparkle = document.createElement('div');
                    sparkle.className = 'sparkle';
                    sparkle.textContent = ['✨', '⭐', '💫', '🌟'][Math.floor(Math.random() * 4)];
                    sparkle.style.left = `${Math.random() * 100}%`;
                    sparkle.style.top = `${Math.random() * 100}%`;
                    sparkle.style.animationDelay = `${Math.random() * 2}s`;
                    
                    sparkleContainer.appendChild(sparkle);
                }, i * 200);
            }
            
            // Remove sparkle container after duration
            setTimeout(() => {
                if (sparkleContainer.parentNode) {
                    sparkleContainer.parentNode.removeChild(sparkleContainer);
                }
            }, duration);
        }
    }
    
    // Enhanced reel sound with pitch variation
    playReelStopSoundWithPitch(reelIndex) {
        if (!this.context) return;
        
        const oscillator = this.context.createOscillator();
        const gainNode = this.context.createGain();
        const filter = this.context.createBiquadFilter();
        
        oscillator.connect(filter);
        filter.connect(gainNode);
        gainNode.connect(this.context.destination);
        
        // Vary pitch based on reel position (left to right gets higher)
        const basePitch = 200 + (reelIndex * 50);
        oscillator.type = 'square';
        oscillator.frequency.setValueAtTime(basePitch, this.context.currentTime);
        oscillator.frequency.exponentialRampToValueAtTime(basePitch * 0.5, this.context.currentTime + 0.1);
        
        filter.type = 'highpass';
        filter.frequency.setValueAtTime(100 + (reelIndex * 20), this.context.currentTime);
        
        gainNode.gain.setValueAtTime(0, this.context.currentTime);
        gainNode.gain.linearRampToValueAtTime(this.volume * 0.4, this.context.currentTime + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.01, this.context.currentTime + 0.15);
        
        oscillator.start(this.context.currentTime);
        oscillator.stop(this.context.currentTime + 0.15);
    }

    // Audio permission prompt system
    setupAudioPermissionPrompt() {
        // Create permission prompt elements
        this.createPermissionPromptHTML();
        
        // Set up first user interaction detection
        this.setupFirstInteractionDetection();
    }
    
    createPermissionPromptHTML() {
        // Create the prompt modal HTML
        const promptHTML = `
            <div id="audio-permission-modal" class="audio-permission-modal" style="display: none;">
                <div class="audio-permission-content">
                    <div class="audio-permission-icon">🔊</div>
                    <h3>Enable Casino Sound Effects?</h3>
                    <p>This game features immersive sound effects including spinning reels, winning celebrations, and coin sounds.</p>
                    <p>Enable audio for the full casino experience!</p>
                    <div class="audio-permission-buttons">
                        <button id="enable-audio-btn" class="audio-permission-btn audio-enable">
                            🎵 Enable Sound
                        </button>
                        <button id="disable-audio-btn" class="audio-permission-btn audio-disable">
                            🔇 Play Silently
                        </button>
                    </div>
                    <div class="audio-permission-note">
                        You can change this later using the volume controls
                    </div>
                </div>
            </div>
        `;
        
        // Add to page if not already present
        if (!document.getElementById('audio-permission-modal')) {
            document.body.insertAdjacentHTML('beforeend', promptHTML);
            this.attachPermissionPromptEvents();
        }
    }
    
    attachPermissionPromptEvents() {
        const enableBtn = document.getElementById('enable-audio-btn');
        const disableBtn = document.getElementById('disable-audio-btn');
        
        if (enableBtn) {
            enableBtn.addEventListener('click', () => {
                this.enableAudioWithPermission();
            });
        }
        
        if (disableBtn) {
            disableBtn.addEventListener('click', () => {
                this.disableAudioWithChoice();
            });
        }
    }
    
    setupFirstInteractionDetection() {
        const handleFirstInteraction = () => {
            // Try to resume audio context immediately on first interaction
            if (this.context && this.context.state === 'suspended') {
                this.context.resume().then(() => {
                    this.audioPermissionGranted = true;
                });
            } else {
                this.audioPermissionGranted = true;
            }
            
            // Remove listeners after first interaction
            document.removeEventListener('click', handleFirstInteraction);
            document.removeEventListener('touchstart', handleFirstInteraction);
            document.removeEventListener('keydown', handleFirstInteraction);
        };
        
        document.addEventListener('click', handleFirstInteraction);
        document.addEventListener('touchstart', handleFirstInteraction);
        document.addEventListener('keydown', handleFirstInteraction);
    }
    
    showAudioPermissionPrompt() {
        if (this.permissionPromptShown) return;
        
        this.permissionPromptShown = true;
        const modal = document.getElementById('audio-permission-modal');
        if (modal) {
            modal.style.display = 'flex';
            
            // Focus management for accessibility
            const enableBtn = document.getElementById('enable-audio-btn');
            if (enableBtn) {
                enableBtn.focus();
            }
            
            // Show a subtle UI indicator that audio is waiting
            this.showAudioPendingIndicator();
        }
    }
    
    hideAudioPermissionPrompt() {
        const modal = document.getElementById('audio-permission-modal');
        if (modal) {
            modal.style.display = 'none';
        }
        this.hideAudioPendingIndicator();
    }
    
    showAudioPendingIndicator() {
        const audioButton = document.getElementById('audio-toggle');
        if (audioButton) {
            audioButton.classList.add('pending-permission');
            audioButton.title = 'Audio permission needed - click to enable';
        }
    }
    
    hideAudioPendingIndicator() {
        const audioButton = document.getElementById('audio-toggle');
        if (audioButton) {
            audioButton.classList.remove('pending-permission');
            audioButton.title = 'Toggle Sound';
        }
    }
    
    enableAudioWithPermission() {
        this.audioPermissionGranted = true;
        this.setEnabled(true);
        
        // Resume audio context
        if (this.context && this.context.state === 'suspended') {
            this.context.resume();
        }
        
        this.hideAudioPermissionPrompt();
        
        // Play a welcome sound
        setTimeout(() => {
            this.playButtonClickSound();
        }, 100);
    }
    
    disableAudioWithChoice() {
        this.audioPermissionGranted = false;
        this.setEnabled(false);
        this.hideAudioPermissionPrompt();
    }

    // Initialize audio context on user interaction (required by modern browsers)
    static initializeOnInteraction() {
        const initAudio = () => {
            if (window.audioManager && window.audioManager.context) {
                window.audioManager.context.resume();
            }
            document.removeEventListener('click', initAudio);
            document.removeEventListener('touchstart', initAudio);
        };
        
        document.addEventListener('click', initAudio);
        document.addEventListener('touchstart', initAudio);
    }
}

// Auto-initialize audio on user interaction
AudioManager.initializeOnInteraction();

// Export for use in other modules
window.AudioManager = AudioManager;
