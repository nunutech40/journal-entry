/**
 * Journal Entry — App JS
 * Alpine.js initialization + HTMX configuration
 */

// HTMX configuration
document.addEventListener('DOMContentLoaded', function() {
    // Configure HTMX defaults
    document.body.addEventListener('htmx:configRequest', function(evt) {
        // Add any default headers here if needed
    });

    // Handle HTMX errors globally
    document.body.addEventListener('htmx:responseError', function(evt) {
        console.error('HTMX error:', evt.detail);
    });

    // Handle HTMX after swap — useful for re-initializing things
    document.body.addEventListener('htmx:afterSwap', function(evt) {
        // Scroll to top of swapped content if it's the main content area
        if (evt.detail.target.classList.contains('l-content')) {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }
    });
});
