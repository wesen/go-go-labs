// visitUrl: https://docs.anthropic.com/en/api/messages
async function clickButtonsAndLogBodyText(rootId) {
  // 1. Find the sub-div by id
  const subDiv = document.getElementById(rootId);
  if (!subDiv) {
    console.error(`No element found with id="${rootId}"`);
    return;
  }
  console.log("subDiv", subDiv);

  // 2. Climb up to the wrapping container (adjust as needed)
  //    Here we go two levels up to that outer .py-6 border div.
  const container = subDiv.parentElement.parentElement;
  if (!container) {
    console.error(`Cannot find container for #${rootId}`);
    return;
  }
  console.log("container", container);

  // 3. Grab all buttons under that container
  const buttons = Array.from(container.querySelectorAll("button"));
  console.log(buttons);

  // 4. Iterate, click, wait for next paint, then log body text
  for (const btn of buttons) {
    console.log("clicking button", btn?.innerText);
    btn.click();

    // wait for the next repaint so any dynamic content has a chance to update
    await new Promise((resolve) => {
      console.log("waiting for next repaint");
      requestAnimationFrame(() => requestAnimationFrame(resolve));
    });

    // Print the text content of the first parent div's next sibling
    const parentDiv = btn.closest("div");
    console.log(parentDiv?.nextElementSibling?.innerText);
  }
}

async function clickAllShowChildAttributes() {
  let totalClicked = 0;

  while (true) {
    // Find all buttons containing "Show child attributes" text
    const buttons = Array.from(document.querySelectorAll("button")).filter(
      (btn) => btn.textContent.trim() === "Show child attributes"
    );

    console.log(
      `Found ${buttons.length} 'Show child attributes' buttons in this iteration`
    );

    if (buttons.length === 0) {
      console.log(
        `Finished expanding all child attributes. Total buttons clicked: ${totalClicked}`
      );
      break;
    }

    // Click each button and wait for repaint
    for (const btn of buttons) {
      console.log("Clicking 'Show child attributes' button");
      btn.click();
      totalClicked++;

      // Wait for next repaint to allow content to update
      await new Promise((resolve) => {
        requestAnimationFrame(() => requestAnimationFrame(resolve));
      });
    }

    // Add a small delay between iterations to ensure DOM is fully updated
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
}

console.clear();
// First click all show child attributes buttons recursively
clickAllShowChildAttributes();
// Then process the specific section
// clickButtonsAndLogBodyText("response-content-citations");

123;
