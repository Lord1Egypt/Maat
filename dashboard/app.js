document.addEventListener('DOMContentLoaded', () => {
    // Inputs & Elements
    const usdInput = document.getElementById('usd-input');
    const egpOutput = document.getElementById('egp-output');
    const volSlider = document.getElementById('volume-slider');
    const volSliderVal = document.getElementById('volume-slider-val');

    // Display fields
    const quoteMid = document.getElementById('quote-mid');
    const quoteSpread = document.getElementById('quote-spread');
    const quoteBuy = document.getElementById('quote-buy');
    const quoteSell = document.getElementById('quote-sell');

    const backingRatioText = document.getElementById('backing-ratio-text');
    const backingRatioFill = document.getElementById('backing-ratio-fill');
    const reserveAssetsVal = document.getElementById('reserve-assets-val');
    const wrappedSupplyVal = document.getElementById('wrapped-supply-val');
    const spreadRevenueVal = document.getElementById('spread-revenue-val');

    // Constants
    const baseMid = 3000.0; // Base mid price of wETH in USD
    const usdToUmaatPeg = 48.25; // 1 USD = 48.25 UMAAT/EGP peg

    // State
    let currentMid = baseMid;
    let accumulatedRevenue = 135892.0;
    let wrappedSupply = 1000000.0;
    let backingRatio = 104.5;

    // Remittance Peg logic
    function updateRemittance() {
        const usdVal = parseFloat(usdInput.value) || 0;
        const rawPeg = usdVal * usdToUmaatPeg;
        // Apply 0.15% spread fee
        const received = rawPeg * (1.0 - 0.0015);
        egpOutput.textContent = received.toLocaleString('en-US', {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }) + " UMAAT";
    }

    // Dynamic Quote Spread Curve logic
    function updateSpreadCurve() {
        const tradeSize = parseInt(volSlider.value);
        volSliderVal.textContent = `${tradeSize} ETH`;

        // Spread base 15bps + 0.5bps per 1 ETH trade size (volatility skew representation)
        const baseBps = 15.0;
        const volMult = 0.5;
        const effSpreadBps = Math.min(100.0, baseBps + tradeSize * volMult); // Caps at 100 bps (1.0%)

        quoteSpread.textContent = `${effSpreadBps.toFixed(1)} bps (${(effSpreadBps / 100).toFixed(2)}%)`;

        const buyPrice = currentMid - currentMid * (effSpreadBps / 10000);
        const sellPrice = currentMid + currentMid * (effSpreadBps / 10000);

        quoteMid.textContent = `$${currentMid.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
        quoteBuy.textContent = `$${buyPrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
        quoteSell.textContent = `$${sellPrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
    }

    // Periodical simulation update
    function simulateMarketFluctuations() {
        // Random minor mid price movement (within 0.1%)
        const pctChange = (Math.random() - 0.5) * 0.002;
        currentMid = currentMid * (1.0 + pctChange);

        // Tick revenue upwards on random trades
        if (Math.random() > 0.4) {
            const addedRev = Math.random() * 8.5;
            accumulatedRevenue += addedRev;
            spreadRevenueVal.textContent = `$${Math.floor(accumulatedRevenue).toLocaleString()}`;
        }

        // Periodically tick wrapped supply slightly
        if (Math.random() > 0.7) {
            const deltaSupply = (Math.random() - 0.45) * 20.0;
            wrappedSupply += deltaSupply;
            wrappedSupplyVal.textContent = `${Math.floor(wrappedSupply).toLocaleString()} wETH`;
        }

        // Recalculate backing ratio
        const assetUnitsValue = wrappedSupply * currentMid * (backingRatio / 100);
        reserveAssetsVal.textContent = `$${Math.floor(assetUnitsValue).toLocaleString()}`;

        updateValidators();
        updateSpreadCurve();
    }

    function updateValidators() {
        const feederList = document.getElementById('feeder-consensus-list');
        if (!feederList) return;
        
        const items = feederList.getElementsByTagName('li');
        for (let i = 0; i < items.length; i++) {
            const dev = (Math.random() - 0.5) * 0.003;
            const price = currentMid * (1.0 + dev);
            const priceSpan = items[i].querySelector('.feeder-price');
            if (priceSpan) {
                priceSpan.textContent = `$${price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
            }
        }
    }

    // Event listeners
    usdInput.addEventListener('input', updateRemittance);
    volSlider.addEventListener('input', updateSpreadCurve);

    // Initial run
    updateRemittance();
    updateSpreadCurve();

    // Start simulation ticks (every 3 seconds)
    setInterval(simulateMarketFluctuations, 3000);
});
