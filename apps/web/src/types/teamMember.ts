import type { SocialMedia } from '@/types/socialMedia';

export type TeamMemberMediaLink = {
  media: SocialMedia;
  link: string;
};

export interface TeamMember {
  id: number;
  name: string;
  role: string;
  isFounder?: boolean;
  socials: TeamMemberMediaLink[];
}
