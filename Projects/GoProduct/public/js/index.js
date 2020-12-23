(function ($) {
    'use strict';
    $(function () {
        var addproduct = function () {
            var $productbox = $('.product_box');
            $productbox.append("<div class='product_box'><p>product</p></div>");
        }
        addproduct();
    });
})(jQuery);