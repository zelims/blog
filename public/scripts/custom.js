// Show dropdown on hover
$('.navbar a[data-toggle="dropdown"]').on("mouseover", function() {
    $(this).next('div.dropdown-menu').addClass('show');
});
$('.navbar .dropdown-menu').on("mouseleave", function() {
    $(this).removeClass('show');
});

$('a.page-number').click(function(e) {
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

$(function() {
    if($("#postSearch")) {
        let searchTimeout = null;
        $("#postSearch").on('keyup', function () {
            $("#postList").html('<div class="text-center"><div class="spinner-border" role="status"><span class="sr-only">Loading...</span></div></div>');
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(function() {
                $.ajax({
                    type: 'POST',
                    url: "/search",
                    headers: { 'REQ_TYPE': "SRV_CALL" },
                    data: $("#searchForm").serialize(),
                    success: function (data) {
                        $("#postList").html(data);
                    },
                    failure: function () {
                        $("#postList").html("Cannot find posts for " + this.val())
                    }
                });
            }, 500);
        });
    }
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
            $('html, body').animate({
                scrollTop: $('.page-container').offset().top
            }, 1000)
        },
        error: function(err, t, s) {
            // TODO: BS4 Modal
            alert("Failed to get data: " + err.responseText);
        }
    });
}

$(function () {
    $('[data-toggle="tooltip"]').tooltip();
});

let _markers = [];
function createMap(element, markers="", seriesdata="") {
    if(markers !== "")
        _markers = markers;
    element.vectorMap({
        map: 'world_mill',
        scaleColors: ['#C8EEFF', '#0071A4'],
        normalizeFunction: 'polynomial',
        zoomButtons: false,
        zoomOnScroll: false,
        panOnDrag: false,
        hoverOpacity: 0.7,
        hoverColor: false,
        markerStyle: {
            initial: {
                fill: '#7E57C2',
                stroke: '#383f47',
                r: 6.5,
            }
        },
        series: {
            regions: [{
                values: seriesdata,
                scale: ['#bf9b9a', '#bf4f50'],
                normalizeFunction: 'polynomial'
            }]
        },
        onRegionTipShow: function(e, el, index){
            if(seriesdata[index] != null)
                el.html(el.html() + ' - ' + seriesdata[index]);
            else
                el.html(el.html())
        },
        backgroundColor: 'transparent',
        regionStyle: {
            initial: {
                fill: '#8d8d8d'
            }
        },
        //markers: _markers,
        onRegionOver: function(e, tip, code) {
            e.preventDefault();
        }
    });

    let serverMapSvg = element.next('.jvectormap-container svg');
    serverMapSvg.removeAttr('width');
    serverMapSvg.removeAttr('height');
}

function showToast(title, type="", time="",  body, duration=2000) {
    let toastTmpl = $('#toast-tmpl');
    let toastId = makeRandomId();
    let toast = toastTmpl.clone().prop('id', toastId);
    toast.removeClass("d-none");

    toast.find("[toast-data=\"toast-header-title\"]").html(title);
    let icon = "far fa-question-circle";
    let toastType = "toast-";

    if (type === "success")
        icon = "fas fa-check";
    else if (type === "error")
        icon = "fas fa-times";
    else if (type === "warning")
        icon = "fas fa-exclamation";
    else
        type = "none";

    toast.addClass(toastType + type);

    toast.find("[toast-data=\"toast-header-icon\"]").addClass(icon);
    toast.find("[toast-data=\"toast-header-time\"]").html(time);
    toast.find("[toast-data=\"toast-body\"]").html(body);

    let shouldHide = true;
    if (duration === -1) {
        shouldHide = false
    }

    $("#toast-list").append(toast);

    toast.toast({
        delay: duration,
        autohide: shouldHide,
    }).toast('show');
}

function makeRandomId() {
    let text = "";
    let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (let i = 0; i < 10; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}