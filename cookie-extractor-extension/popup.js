let extractedCookies = '';

document.getElementById('extractBtn').addEventListener('click', extractCookies);
document.getElementById('copyBtn').addEventListener('click', copyToClipboard);

async function extractCookies() {
  const statusDiv = document.getElementById('status');
  const resultDiv = document.getElementById('result');
  const copyBtn = document.getElementById('copyBtn');
  
  try {
    statusDiv.style.display = 'block';
    statusDiv.className = 'info';
    statusDiv.textContent = 'Extracting cookies...';
    
    // Get all cookies from icloud.com domain
    const cookies = await chrome.cookies.getAll({
      domain: '.icloud.com'
    });
    
    if (cookies.length === 0) {
      statusDiv.className = 'error';
      statusDiv.textContent = '❌ No cookies found! Make sure you are logged in to icloud.com';
      resultDiv.style.display = 'none';
      copyBtn.disabled = true;
      return;
    }
    
    // Format as "name=value; name=value; ..."
    extractedCookies = cookies
      .map(cookie => `${cookie.name}=${cookie.value}`)
      .join('; ');
    
    // Display result
    resultDiv.textContent = extractedCookies;
    resultDiv.style.display = 'block';
    
    // Enable copy button
    copyBtn.disabled = false;
    
    // Check for required cookies
    const hasWebAuthUser = cookies.some(c => c.name === 'X-APPLE-WEBAUTH-USER');
    const hasWebAuthToken = cookies.some(c => c.name === 'X-APPLE-WEBAUTH-TOKEN');
    const hasDsWebSession = cookies.some(c => c.name === 'X-APPLE-DS-WEB-SESSION-TOKEN');
    
    if (hasWebAuthUser || hasDsWebSession) {
      statusDiv.className = 'success';
      statusDiv.textContent = `✓ Found ${cookies.length} cookies (including HttpOnly)`;
    } else {
      statusDiv.className = 'error';
      statusDiv.textContent = `⚠️ Found ${cookies.length} cookies but missing critical auth cookies. Please log in to icloud.com first.`;
    }
    
  } catch (error) {
    statusDiv.className = 'error';
    statusDiv.textContent = '❌ Error: ' + error.message;
    resultDiv.style.display = 'none';
    copyBtn.disabled = true;
  }
}

async function copyToClipboard() {
  const statusDiv = document.getElementById('status');
  
  try {
    await navigator.clipboard.writeText(extractedCookies);
    statusDiv.className = 'success';
    statusDiv.textContent = '✓ Copied to clipboard! Paste into cookies.txt';
  } catch (error) {
    statusDiv.className = 'error';
    statusDiv.textContent = '❌ Failed to copy: ' + error.message;
  }
}
