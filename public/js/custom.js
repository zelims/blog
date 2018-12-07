$('a.page-number').click(function(e)
{
    let page = $(this).attr('page-number');

    switchPage(page);

    return false;
});

$('li#nextPage a').click(function(e) {
    let nextPage = +getCurrentPage() + 1;
    switchPage(nextPage);
});
$('li#prevPage a').click(function(e) {
    let prevPage = +getCurrentPage() - 1;
    switchPage(prevPage);
});

function getPageCount() {
    return $('#postPageCount').val();
}

function getCurrentPage() {
    return $('#curPage').val();
}

function switchPage(page) {
    $('li.page-item.active').removeClass('active');
    $('[page-number="'+page+'"]').parent().addClass('active');
    let pageCount = getPageCount();
    let prevPageBtn = $('#prevPage');
    let nextPageBtn = $('#nextPage');

    if(page == 1) {
        prevPageBtn.addClass('disabled');
    } else {
        prevPageBtn.removeClass('disabled');
    }
    if(page == pageCount) {
        nextPageBtn.addClass('disabled');
    } else {
        nextPageBtn.removeClass('disabled');
    }
    $.ajax({
        url: "/page/" + page,
        method: "post",
        success: function(data) {
            $('#curPage').val(page);
            $('#post-list').html(data);
        },
        error: function(err, t, s) {
            alert("Failed to get data: " + err.responseText);
        }
    });
}