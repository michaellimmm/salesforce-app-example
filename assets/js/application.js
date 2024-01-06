"use strict";

async function copyToClipboard(element) {
    try {
        await navigator.clipboard.writeText(element.value);
        alert("Text copied to clipboard")
    } catch (error) {
        alert("Failed to copy to clipboard")
    }
}
