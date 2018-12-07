$('a.page-link').click(function(e)
{
    $('li.page-item.active').removeClass('active');
    $(this).parent().addClass('active');
    let page = $(this).attr('page-number');
    $.ajax({
        url: "/page/" + page,
        method: "post",
        success: function(data) {
            $('#post-list').html(data);
        },
        error: function(err, t, s) {
            alert("Failed to get data: " + err.responseText);
        }
    });
    return false;
});