
function gallery(elementSelector, pathPrefix, minContainerWidth, maxContainerCount, getOptions) {

    var imagePath = pathPrefix + '/image/';
    var images = [];
    if ('ontouchstart' in window || 'onmsgesturechange' in window) {

    }
    
    $(elementSelector).after(
        '<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">'
            + '  <div class="pswp__bg"></div>'
            + '  <div class="pswp__scroll-wrap">'
            + '    <div class="pswp__container"><div class="pswp__item"></div><div class="pswp__item"></div><div class="pswp__item"></div></div>'
            + '    <div class="pswp__ui pswp__ui--hidden">'
            + '      <div class="pswp__top-bar">'
            + '        <div class="pswp__counter"></div>'
            + '        <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>'
            + '        <button class="pswp__button pswp__button--share" title="Share"></button>'
            + '        <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>'
            + '        <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>'
            + '        <div class="pswp__preloader"><div class="pswp__preloader__icn"><div class="pswp__preloader__cut"><div class="pswp__preloader__donut"></div></div></div></div>'
            + '      </div>'
            + '      <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap"><div class="pswp__share-tooltip"></div></div>'
            + '      <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)"></button>'
            + '      <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)"></button>'
            + '      <div class="pswp__caption"><div class="pswp__caption__center"></div></div>'
            + '    </div>'
            + '  </div>'  
            + '</div>');

    var clickHandler = function(event) {
        event.preventDefault();
        var imageIndex = $(this).attr("data-gindex");
        if (imageIndex) {
            imageIndex = parseInt(imageIndex);
        } else {
            imageIndex = 0;
        }
        var pswpElement = document.querySelectorAll('.pswp')[0];
        var photoSwipe;
        var options = {
            index: imageIndex,
            shareEl: true,
            shareButtons: [
                {id:'facebook', label:'Teilen auf Facebook', url:'https://www.facebook.com/sharer/sharer.php?u={{url}}'},
                {id:'pinterest', label:'Pin it', url:'http://www.pinterest.com/pin/create/button/?url={{url}}&media={{image_url}}&description={{text}}'},
            ],
            getPageURLForShare: function( shareButtonData ) {
                var image = photoSwipe.currItem;
                return location.protocol + '//' + location.hostname + "/share" + image.src
                    +'?title=' + encodeURIComponent('by '+image.user.nickName);
            },
            getTextForShare: function( shareButtonData ) {
                var image = photoSwipe.currItem;
                return 'by '+image.user.nickName;
            }
        };
        photoSwipe = new PhotoSwipe(pswpElement, PhotoSwipeUI_Default, images, options);
        photoSwipe.init();
    }

    var createFacebookClickHandler = function(image) {
        return function(event) {            
            event.preventDefault();
            var shareUrl = encodeURIComponent(location.protocol + '//' + location.hostname + "/share" + image.src
                                              +'?title=' + encodeURIComponent('by '+image.user.nickName));
            window.open('http://www.facebook.com/sharer/sharer.php?u='+shareUrl, 'facebook_share', 'height=320, width=640, toolbar=no, menubar=no, scrollbars=no, resizable=no, location=no, directories=no, status=no');
        }
    }

    var containerCount;
    
    var layout = function() {
        var element = $(elementSelector);
        var width = element.width();
        var newContainerCount = Math.min(Math.max(1, Math.floor(width/minContainerWidth)), maxContainerCount);
        if (newContainerCount == containerCount) {
            return;
        }
        containerCount = newContainerCount        
        var containerWidthPercent = 100*(width/containerCount)/width;

        element.empty();
        
        var container = [];
        for (var c=0; c<containerCount; c++) {
            container.push($('<div class="image-list" style="width: '+containerWidthPercent+'%;"/>').appendTo(element));
        }
        
        for (var i=0; i<images.length; i++) {
            var smallestContainer = container[0];
            for (var c=1; c<container.length; c++) {
                if (container[c].height() < smallestContainer.height()) {
                    smallestContainer = container[c];
                }
            }                      

            var div = $('<div class="gallery-thumbnail"/>')
                .appendTo(smallestContainer);

            $('<a href="#" class="share"/>')
                .click(createFacebookClickHandler(images[i]))
                .appendTo(div);
            
            $('<a href="#" class="gallery-thumbnail-link" data-gindex="'+i+'"><img src="' + images[i].msrc + '"></a>')
                .click(clickHandler)
                .appendTo(div);

            var href ='';
            var link = images[i].user.link;
            if (link) {
                link = link.trim();
                if (link.indexOf('http') != 0) {
                    link = 'http://' + link;
                }
                var href = ' href="'+link+'"';
            }
            $('<div class="gallery-thumbnail-description"><a target="_blank"'+href+'>by '+images[i].user.nickName+'</a></div>')
                .appendTo(div);
        }
    };

    var fetch = function() {
        $.getJSON(pathPrefix + "/api/images?"+ getOptions, function( data ) {
            images = data;
            $.each( images, function( key, image ) {
                image.src = imagePath + image.src
                image.msrc = imagePath + image.msrc
            });

            layout();
        });
    }

    $( window ).resize(layout);
    fetch();
}
