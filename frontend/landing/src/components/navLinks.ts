type NavLink = { href: string; label: string };

export const navLinks: NavLink[] = [
	{
		href: '#features',
		label: 'Features',
	},
	{
		href: '#integrations',
		label: 'Integrations',
	},
	{
		href: '#reviews',
		label: 'Reviews',
	},
	{
		href: '#team',
		label: 'Team',
	},
];

export const scrollToSection = (target: HTMLElement) => {
	const sectionY = window.scrollY - target.getBoundingClientRect().top;

	window.scrollTo({
		top: sectionY,
		behavior: 'smooth',
	});
};
