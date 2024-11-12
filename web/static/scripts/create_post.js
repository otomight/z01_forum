document.getElementById("createPostForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the default form submission
    
    const title = document.getElementById("title").value.trim();
    const content = document.getElementById("content").value.trim();
    const category = document.getElementById("category").value.trim();
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

    //Create URLSearchParams to encode form data 
    const formData = new URLSearchParams();
    formData.append("title", title);
    formData.append("content", content);
    formData.append("category", category);
    formData.append("tags", tags);
    
    // submit Formdata to server 
    fetch("/api/posts/create", {
        method: "POST",
        body: formData
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


