/**
 * @typedef {number}
 */
var initialInputOlder = 0; // miliseconds Unix TimeStamp
/**
 * @typedef {number}
 */
var initialInputNewer = 0; // miliseconds Unix TimeStamp


/**
 * @typedef {HTMLInputElement}
 */
var olderThanInputElement;
/**
 * @typedef {HTMLInputElement}
 */
var newerThanInputElement;

document.addEventListener(`DOMContentLoaded`, initializePage)

function initializePage() {

    // Update the page's time filter
    initialInputOlder = Number(document.getElementById(`inp-older`).textContent.trim())
    initialInputNewer = Number(document.getElementById(`inp-newer`).textContent.trim())

    // Bind the date elements
    olderThanInputElement = document.getElementById(`olderThan`);
    newerThanInputElement = document.getElementById(`newerThan`);

    olderThanInputElement.value = convertTimestampToDateTimeLocal(initialInputOlder);
    newerThanInputElement.value = convertTimestampToDateTimeLocal(initialInputNewer);

}

// Handle one of the date elements having their value changed.
function applyDateInterval() {
    const olderTimeStamp = new Date(olderThanInputElement.value).getTime()
    const newerTimeStamp = new Date(newerThanInputElement.value).getTime()

    window.location.href = `/?older=${olderTimeStamp}&newer=${newerTimeStamp}`;
}

/**
 * Converts a Unix timestamp to a standard datetime string
 * @param {number} timestamp - The Unix timestamp in milliseconds.
 * @returns {string} - A normal string with Y-m-d H:i:s format
 */
function convertTimestampToDateTimeLocal(timestamp) {
    const date = new Date(timestamp);
    const offset = date.getTimezoneOffset() * 60000; // offset in milliseconds
    const localDate = new Date(date.getTime() - offset);
    return localDate.toISOString().slice(0, 19);
}
