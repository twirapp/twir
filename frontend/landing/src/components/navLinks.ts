export const navLinks: { href: string; label: string }[] = [
  {
    href: '#features',
    label: 'Features',
  },
  {
    href: '#reviews',
    label: 'Reviews',
  },
  {
    href: '#team',
    label: 'Team',
  },
  // {
  //   href: '#',
  //   label: 'Pricing',
  // },
];

export const scrollToSection = (target: HTMLElement) => {
	const sectionY = window.scrollY - target.getBoundingClientRect().top;

  window.scrollTo({
    top: sectionY,
    behavior: 'smooth',
  });
};
