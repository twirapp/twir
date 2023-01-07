import { headerHeightStore } from './header.js';

export const scrollToLandingSection = (target: HTMLElement, smooth = true) => {
  const sectionId = target.dataset.section;
  if (!sectionId) {
    throw new Error(`Target has not data-section attribute`);
  }

  const section = document.getElementById(sectionId);
  if (!section) {
    throw new Error(`Section "${sectionId}" is not found`);
  }

  const sectionY = window.scrollY - headerHeightStore.get() + section.getBoundingClientRect().top;

  window.scrollTo({
    top: sectionY,
    behavior: smooth ? 'smooth' : 'auto',
  });
};
