import * as Prisma from '@tsuwari/prisma';

const { PrismaClient } = Prisma;

export const prisma = new PrismaClient();
