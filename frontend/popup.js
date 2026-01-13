document.getElementById('scrapeBtn').addEventListener('click', async () => {
  let [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  chrome.scripting.executeScript({
    target: { tabId: tab.id },
    func: () => {
      console.log("--- EXECUTING FINAL DYNAMIC SCRAPE & SEND ---");
      
      const url = window.location.href;
      const placeIdMatch = url.match(/place\/([^\/]+)/);
      const locationId = placeIdMatch ? placeIdMatch[1] : "Unknown";

      const cards = document.querySelectorAll('.jftiEf');
      
      const results = Array.from(cards).map(card => {
        // 1. PROFILE PHOTO: Targeted header scan
        const avatarImg = card.querySelector('button img, .NBa79c, img[src*="googleusercontent.com/a/"]');
        const profile_photo = avatarImg ? avatarImg.src : "No Profile Image";

        // 2. MEDIA: Broad scan using your successful "Recovery" logic
        const allElements = Array.from(card.querySelectorAll('img, button[style*="background-image"]'));
        const mediaUrls = allElements
            .map(el => {
                if (el.tagName === 'IMG') return el.src;
                const bg = window.getComputedStyle(el).backgroundImage;
                return bg.match(/url\("?(.+?)"?\)/)?.[1];
            })
            // Filters out placeholders, duplicates, and the profile photo
            .filter(src => src && src !== profile_photo && !src.includes('data:image'));

        // 3. FINAL DATA OBJECT: Matches Go Struct exactly
        return {
          location_id: locationId,
          review_id: card.getAttribute('data-review-id') || "N/A",
          author: card.querySelector('.d4r55')?.innerText || "Unknown",
          rating: card.querySelector('.kvMYJc')?.getAttribute('aria-label')?.match(/\d+/)?.[0] || "0",
          profile_photo: profile_photo,
          text: card.querySelector('.wiI7pd')?.innerText || "No text",
          media_links: mediaUrls.join(', ') // Joined string for SQLite TEXT column
        };
      });

      console.table(results);

      // --- SEND TO BACKEND ---
      fetch('http://localhost:8080/api/reviews', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(results)
      })
      .then(response => {
        if (response.ok) console.log("✅ Data successfully saved to Database.");
      })
      .catch(err => console.error("❌ Connection failed:", err));
    }
  });
});