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
            alert("Failed to get data: " + err.responseText);
        }
    });
}

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
/**
 * Show Notification at top of screen
 * @param parent        (null)              =   parent div (a, btn, etc.)
 * @param target_id     (notify-value)      =   id of the target notify (data-value attr on parent)
 * @param text          (notify-text)       =   text of message
 * @param scheme        (notify-scheme)     =   color scheme (bootstrap)
 * @param timeout       (notify-timeout)    =   timeout in ms
 */
function showNotifyAlert(parent, target_id="", text="", scheme="", timeout=2500) {

    target_id = '#' + target_id;

    $(target_id).addClass('show');

    if(text === "" && parent)
        text = parent.attr('notify-text');

    $(target_id).find('#notify-alert-text').html(text);

    if(scheme === "" && parent)
        scheme = parent.attr('notify-scheme');

    if(timeout === 2500 && parent) {
        let dto = parent.attr('notify-timeout');
        if(dto)
            timeout = dto;
    }

    var icon;
    switch(scheme) {
        case "danger":
            icon = "fa-times";
            break;
        case "warning":
            icon = "fa-exclamation-circle";
            break;
        case "success":
            icon = "fa-check";
            break;
        case "info":
            icon = "fa-question";
            break;
        default:
            scheme = "primary";
            icon = "fa-exclamation";
    }
    $(target_id).find('#notify-alert-scheme').addClass('bg-' + scheme);
    $(target_id).find('#notify-alert-icon').addClass(icon);
    setTimeout(function() {
        $(target_id).removeClass("show");
    }, timeout);
}