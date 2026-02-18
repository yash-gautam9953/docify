// ─── Copy to clipboard ──────────────────────────────────────────
function copyCode(btn) {
  const code = btn.nextElementSibling.textContent;
  navigator.clipboard.writeText(code).then(() => {
    btn.textContent = "✅ Copied!";
    setTimeout(() => {
      btn.textContent = "📋 Copy";
    }, 2000);
  });
}

// ─── Scroll fade-in animation ───────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
  // Add fade-in class to all animated elements
  const targets = document.querySelectorAll(
    ".card, .step, .install-card, .detect-item, .port-priority, .compare-col, .release-card, .terminal",
  );
  targets.forEach((el) => el.classList.add("fade-in"));

  // Intersection Observer for scroll animations
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry, index) => {
        if (entry.isIntersecting) {
          // Stagger animation
          setTimeout(() => {
            entry.target.classList.add("visible");
          }, index * 80);
          observer.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.15, rootMargin: "0px 0px -40px 0px" },
  );

  targets.forEach((el) => observer.observe(el));
});
