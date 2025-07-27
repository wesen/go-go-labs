# Casino Slot Machine Game

A fully responsive, feature-rich casino slot machine game built with vanilla HTML, CSS, and JavaScript.

## Features

### 🎰 Core Gameplay
- **3-reel slot machine** with animated spinning
- **Multiple symbols** with different payout values
- **Progressive betting** system ($1 - $100)
- **Real-time balance** tracking with persistence
- **Win detection** and payout calculation
- **Game history** tracking

### 🎨 Visual Design
- **Casino-themed UI** with golden accents and gradients
- **Smooth animations** for spinning reels and wins
- **Visual feedback** for wins, losses, and jackpots
- **Responsive design** that works on all devices
- **CSS3 effects** including glow, pulse, and celebration animations

### 🔊 Audio System
- **Web Audio API** powered sound engine
- **Dynamic sound effects** for spins, wins, and interactions
- **Layered audio** for different win levels
- **Volume controls** and mute functionality
- **Synthesized sounds** (no external audio files needed)
- **Autoplay policy handling** with permission prompt
- **Cross-browser audio compatibility** with fallbacks

### 📱 Responsive Design
- **Mobile-first** approach with touch-friendly controls
- **Tablet optimization** with reorganized layout
- **Desktop enhancements** with larger displays
- **Landscape mode** support for mobile devices
- **Print-friendly** styles

### 🎮 User Experience
- **Keyboard shortcuts** for power users
- **Visual feedback** for all interactions
- **Error handling** with graceful degradation
- **Performance monitoring** and optimization
- **Local storage** for game state persistence
- **Button debouncing** to prevent rapid clicking issues
- **WCAG-compliant accessibility** with screen reader support

## Quick Start

1. **Open the game**: Simply open `index.html` in a web browser
2. **Set your bet**: Use the input field or preset buttons
3. **Spin the reels**: Click the red SPIN button or press Space/Enter
4. **Check results**: Watch for winning combinations and collect payouts

## Controls

### Mouse/Touch
- **Spin Button**: Start a new spin
- **Bet Input**: Manually enter bet amount
- **Preset Buttons**: Quick bet amount selection
- **Double-click Balance**: Add debug money (development only)

### Keyboard Shortcuts
- **Space** or **Enter**: Spin the reels
- **1-5**: Set bet to preset amounts ($1, $5, $10, $25, $100)
- **↑/↓**: Increase/decrease bet by $1
- **M**: Toggle audio on/off

## Paytable

| Symbols | Payout Multiplier |
|---------|-------------------|
| 🍀🍀🍀  | x200             |
| 💎💎💎  | x100             |
| ⭐⭐⭐  | x75              |
| 🔔🔔🔔  | x50              |
| 🍇🍇🍇  | x30              |
| 🍒🍒🍒  | x20              |
| 🍋🍋🍋  | x15              |
| 🍊🍊🍊  | x10              |

*Two matching symbols also pay smaller amounts*

## Technical Architecture

### Component Structure
```
CasinoSlotMachine (main.js)
├── GameState (game-state.js)      - Balance, betting, persistence
├── SlotMachine (slot-machine.js)  - Reel animation, game logic
├── AudioManager (audio-manager.js) - Sound effects, Web Audio API
└── UIController (ui-controller.js) - User interface, interactions
```

### File Organization
```
casino-slot-machine/
├── index.html              - Main game page
├── css/
│   ├── styles.css          - Base styles and layout
│   ├── slot-machine.css    - Slot machine specific styles
│   └── responsive.css      - Mobile and responsive styles
├── js/
│   ├── game-state.js       - Game state management
│   ├── slot-machine.js     - Slot machine logic
│   ├── audio-manager.js    - Audio system
│   ├── ui-controller.js    - UI interactions
│   └── main.js             - Application initialization
├── assets/
│   ├── images/            - Game images (empty - uses emoji)
│   └── sounds/            - Sound files (empty - uses Web Audio)
├── tests/                 - Test files
└── README.md              - This file
```

## Browser Compatibility

### Minimum Requirements
- **Modern browsers** with ES6+ support
- **Web Audio API** for sound (optional)
- **Local Storage** for game persistence (optional)

### Tested Browsers
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+
- Mobile Safari (iOS 13+)
- Chrome Mobile (Android 8+)

## Performance Features

- **RequestAnimationFrame** for smooth animations
- **CSS Hardware Acceleration** for transforms
- **Efficient DOM updates** with minimal reflows
- **Memory management** for audio contexts
- **Lazy loading** of non-critical components

## Development

### Debug Mode
Add `?debug=true` to the URL or run on localhost to enable:
- **Console access** to all game components
- **Debug balance** addition (double-click balance)
- **Enhanced logging** and error reporting

### Customization
The game is highly modular and customizable:
- **Symbols**: Modify `gameState.symbols` array
- **Payouts**: Update `gameState.payouts` object
- **Colors**: Edit CSS custom properties
- **Sounds**: Extend `AudioManager` class methods

### Testing
Open browser developer tools and use:
```javascript
// Access game components
window.casino.getDebugInfo()
window.gameState.addDebugBalance(1000)
window.slotMachine.forceResult(['🍀', '🍀', '🍀'])
```

## Security Considerations

- **Client-side only**: All logic runs in browser (not real money)
- **No external dependencies**: Self-contained code
- **No data transmission**: Everything stays local
- **Safe random numbers**: Uses Math.random() for entertainment only

## License

This project is provided as-is for educational and entertainment purposes.

## Credits

Built with vanilla web technologies:
- **HTML5** for structure
- **CSS3** for styling and animations  
- **JavaScript ES6+** for logic
- **Web Audio API** for sound synthesis
- **Emoji symbols** for graphics (universal support)

---

*Have fun and play responsibly! This is a demo game for entertainment only.*
