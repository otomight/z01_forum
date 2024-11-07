document.getElementById("createPostForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the default form submission
    
    const title = document.getElementById("title").value.trim();
    const content = document.getElementById("content").value.trim();
    const tags = document.getElementById("tags").value.trim();
    const errorMessage = document.getElementById("error-message");
    const successMessage = document.getElementById("success-message");

    // Reset messages
    errorMessage.textContent = "";
    successMessage.textContent = "";

    // Basic validation
    if (!title || !content) {
        errorMessage.textContent = "Title and content are required!";
        return;
    }

    // Optionally format tags
    const formattedTags = tags.split(",").map(tag => tag.trim()).filter(tag => tag).join(", ");
    
    // Simulating a server submission with Fetch API (this would actually be to your backend in a real app)
    fetch("/posts/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ title, content, tags: formattedTags })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            successMessage.textContent = "Post created successfully!";
            document.getElementById("createPostForm").reset(); // Clear form after success
        } else {
            errorMessage.textContent = "There was an issue creating the post.";
        }
    })
    .catch(error => {
        console.error("Error:", error);
        errorMessage.textContent = "Error connecting to the server.";
    });
});

// Optional: Format tags as the user types
document.getElementById("tags").addEventListener("input", function(event) {
    const input = event.target;
    const formattedTags = input.value.split(",").map(tag => tag.trim()).filter(tag => tag).join(", ");
    input.value = formattedTags;
});
