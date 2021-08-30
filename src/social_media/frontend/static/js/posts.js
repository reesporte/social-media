// This file contains functions needed for getting and submitting posts
/**
 * Given a post JSON object, create a post div and return it
 * @param {*} post json object with post data
 * @returns post div
 */
function formatPost(post) {
    // open div
    if ((post.hasOwnProperty('replyingTo')) && (post.replyingTo != "")) {
        var postData = "<div id='" + post.id + "' class='post reply'>";
    } else {
        var postData = "<div id='" + post.id + "' class='post'>";
    }
    postData += "<span id='offset" + post.id + "' style='display:none'>0</span>";

    // post metadata
    postData += '<span class="metadata">';

    // author
    postData += '<span class="author"><a class="profilelink" href="/profile/' + post.username + '">' + post.username + '</a></span>';

    // datetime
    postData += '<span class="datetime">'
    let time = new Date(((post.unixTime) * 1000));
    if ((time.getFullYear() === new Date().getFullYear()) && (time.getMonth() === new Date().getMonth()) && (time.getDate() === new Date().getDate())) {
        postData += '<em>' + time.toLocaleTimeString() + '</em>';
    } else {
        postData += '<em>' + time.toLocaleString() + '</em>';
    }
    postData += '</span>'

    // delete button
    if (post.username === localStorage.getItem("username")) {
        postData += '<button class="button-link delete" onclick=deletePost("' + post.id + '") title="delete your post">delete</button>'
    }
    postData += '</span>'

    // media
    if (post.hasOwnProperty('media') && post.media != "") {
        postData += "<embed autoplay='false' autostart='0' id='media" + post.id + "' src='" + post.media + "' class='media' width=300>"
    }

    // actual post content
    postData += '<span class="message">' + post.message + '</span>'

    // reply button
    postData += '<button class="button-link" onclick=replyToPost("' + post.id + '") title="reply to this post">reply</button>';

    // show replies button
    if ((post.hasOwnProperty('numReplies')) && (post.numReplies > 0)) {
        postData += '<button class="button-link" onClick=showReplies("';
        postData += post.id + '") title="show/hide replies to this post">[<span id="numReplies'+post.id+'">' + post.numReplies + '</span>]</button>'
    }

    // close div
    postData += "</div>"

    return postData;
}

/**
 * get an individual post by id
 * @param {*} id 
 * @returns 
 */
async function getPostByID(id) {
    const response = await fetch('/posts/' + id);
    const json = await response.json();
    return json;
}

/**
 * get up to 10 post replies by post id
 * @param {*} id 
 * @returns 
 */
async function get20PostReplies(id) {
    const json = await fetch('/posts/replies/' + id, {
        headers: {
            'Offset': document.getElementById('offset' + id).innerText
        }
    }).then(function (response) {
        if (response.ok) {
            offset = document.getElementById('offset' + id).innerText;
            offset += 10;
            document.getElementById('offset' + id).innerText = offset;
            return response.json();
        } else {
            return JSON.parse("[]")
        }
    });

    return json;
}


/**
 * show the replies to the post with the given id
 * @param {*} id 
 */
async function showReplies(id) {
    post = document.getElementById(id);
    replies = document.getElementById('replies' + id);
    if (replies === null) {
        var replies = '<div id="replies' + id + '" style="display:block">';
        var json = await get20PostReplies(id);
        for (let i = 0; i < json.length; i++) {
            replies += formatPost(json[i]);
        }
        
        replies += '</div>';
        post.innerHTML = post.innerHTML + replies;
    } else {
        toggleCollapsed('replies' + id);
    }
}

/**
 * load more replies to the post with the given id
 * @param {*} id 
 */
async function loadMoreReplies(id) {
    replies = document.getElementById('replies' + id);
    var json = await get20PostReplies(id);
    for (let i = 0; i < json.length; i++) {
        replies.innerHTML = replies.innerHTML + formatPost(json[i]);
    }
}

/**
 * Given an array of json data containing post information, 
 * add each post to the 'posts' div
 * @param {*} data 
 * @param {boolean} reload whether to reload the page
 */
async function addPosts(data, reload = true) {
    var mainContainer = document.getElementById('posts');
    var postHtml = "";
    for (var i = 0; i < data.length; i++) {
        postHtml += formatPost(data[i]);
    }
    if (reload) {
        document.getElementById('posts-offset').innerText = 10;
        mainContainer.innerHTML = postHtml;
    } else {
        mainContainer.innerHTML += postHtml;
    }
}

/**
 * Load all posts returned from '/posts'
 * into the DOM
 */
async function loadPosts() {
    toggleCollapsed('spinner-container');
    fetch('/posts/')
        .then(function (response) {
            if (response.redirected) {
                window.location.replace("/login/"); 
            } else if (response.ok) {
                return response.json();
            } else if (response.status === 404) {
                throw 'you\'re all caught up! : )'
            } else {
                throw 'couldn\'t load more posts : ('
            }
        })
        .then(function (data) {
            addPosts(data);
        }).catch(function (error) {
            alert(error)
        });
    toggleCollapsed('spinner-container');
}

function loadMorePosts() {
    toggleCollapsed('spinner-container');
    let offset = document.getElementById('posts-offset').innerText;
    fetch('/posts/', {
        headers: {
            'Offset': offset
        }
    }).then(function (response) {
        if (response.ok) {
            return response.json();
        } else if (response.status === 404) {
            throw 'No more posts to load!';
        }
    }).then(function (data) {
        document.getElementById('posts-offset').innerText = parseInt(offset) + 10;
        addPosts(data, false);
    }).catch(function () {
        alert('you\'re all caught up!');
    });
    toggleCollapsed('spinner-container');
}

/**
 * Prep DOM to reply to a post
 * @param {string} replyingTo 
 */
function replyToPost(replyingTo) {
    toggleCollapsed("post-modal", replyingTo);
}

/**
 * Submits user message to /posts endpoint
 */
async function submitPost() {
    await fetch('/posts/', {
        method: 'POST',
        body: JSON.stringify({
            "message": document.getElementById('message').value,
            "replyingTo": document.getElementById('replyingTo').value
        })
    }).then(function (response) {
        if (response.ok) {
            document.getElementById("post-form").reset();
            return response.json();
        } else {
            alert("couldn't send your message : (")
        }
    }).then(function (json) {
        let files = document.getElementById("media").files
        if (files.length > 0 && files[0] != undefined && files[0] != null) {
            document.getElementById("postId").value = json['id'];
            return fetch('/media/', {
                method: 'POST',
                body: new FormData(document.getElementById("media-form"))
            }).then(function (response) {
                if (!response.ok) {
                    alert("could not attach image to post");
                }
                document.getElementById("media-form").reset();
            }).then(function () {
                return true;
            });
        }
        return false;
    }).then(function (mediaPosted) {
        if (mediaPosted) {
            return window.location.reload();
        } else {
            toggleCollapsed("post-modal");
            loadPosts();
        }
    })


}

/**
 * Delete post with given id
 */
function deletePost(id) {
    fetch('/posts/' + id, {
        method: 'DELETE'
    }).then(function (response) {
        console.log(response.status);
        if (response.ok) {
            loadPosts()
        } else {
            alert("couldn't delete post : (")
        }
    });
}
