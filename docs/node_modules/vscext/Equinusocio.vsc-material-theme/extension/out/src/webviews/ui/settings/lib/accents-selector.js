"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const templateSingleAccent = (accentName, accentColor) => {
    const dashAccentName = accentName.toLowerCase().replace(/ /gi, '-');
    return `
    <label for="${dashAccentName}" data-color="${accentColor}">${accentName}</label>
    <input type="radio" name="accents" id="${dashAccentName}" value="${dashAccentName}" />
  `;
};
exports.default = (containerSelector, accentsObject, currentAccent) => {
    const container = document.querySelector(containerSelector);
    for (const accentKey of Object.keys(accentsObject)) {
        const el = document.createElement('div');
        el.innerHTML = templateSingleAccent(accentKey, accentsObject[accentKey]);
        if (accentKey === currentAccent) {
            el.setAttribute('selected', 'true');
            el.querySelector('input').setAttribute('checked', 'checked');
        }
        container.appendChild(el);
    }
};
//# sourceMappingURL=accents-selector.js.map