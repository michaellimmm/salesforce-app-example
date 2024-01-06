"use strict";

function _getCarouselParent() {
    return document.querySelector(".carousel-inner");
}

function _getCarouselItems() {
    const parent = _getCarouselParent()
    return parent.querySelectorAll(".carousel-item");
}

function _getCarouselBtnNext() {
    return document.querySelector(".carousel-btn-next");
}

function _getCarouselBtnPrev() {
    return document.querySelector(".carousel-btn-prev");
}

function showCarouselBtnNext() {
    _showElement(_getCarouselBtnNext());
}

function hideCarouselBtnNext() {
    _hideElement(_getCarouselBtnNext());
}

function showCarouselBtnPrev() {
    _showElement(_getCarouselBtnPrev());
}

function hideCarouselBtnPrev() {
    _hideElement(_getCarouselBtnPrev());
}

function _showElement(element) {
    element.style.visibility = 'visible';
}

function _hideElement(element) {
    element.style.visibility = 'hidden';
}

function hideCarouselBtnPrevOnFirstSlide() {
    const carouselItem = _getCarouselItems();
    if (isCarouselActive(carouselItem[0])) {
        hideCarouselBtnPrev();
    } else {
        showCarouselBtnPrev();
    }
}

function onCarouselBtnNextClicked() {
    const items = _getCarouselItems();
    const currActiveIndex = _getActiveCarouselIndex(items);
    const newIndex = currActiveIndex + 1;
    _renderCarouselBtn(newIndex, items.length);
}

function onCarouselBtnPrevClicked() {
    const items = _getCarouselItems();
    const currActiveIndex = _getActiveCarouselIndex(items);
    const newIndex = currActiveIndex - 1;
    _renderCarouselBtn(newIndex, items.length);
}

function _getActiveCarouselIndex(items) {
    let activeIndex = 0;
    for (let i = 0; i < items.length; i++) {
        if (isCarouselActive(items[i])) {
            activeIndex = i;
            break;
        }
    }

    return activeIndex;
}

function isCarouselActive(carouselItem) {
    return carouselItem.classList.contains("active");
}

function _renderCarouselBtn(currSlideIndex, totalSlide) {
    if (currSlideIndex < totalSlide - 1) {
        showCarouselBtnNext();
    } else {
        hideCarouselBtnNext();
    }

    if (currSlideIndex > 0) {
        showCarouselBtnPrev();
    } else {
        hideCarouselBtnPrev();
    }
}