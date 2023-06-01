import { PrismaClient } from '@prisma/client';
import { faker } from '@faker-js/faker';
// @ts-ignore
import Snowflake from 'snowflake-id';

const prisma = new PrismaClient();
const snowflake = new Snowflake({
  mid: 42,
  offset: (2023 - 1970) * 31536000 * 1000,
});

const USERS_NUMBER = 500;
const ROOMS_NUMBER = 500;

const generateId = () => {
  return BigInt(snowflake.generate());
};

const randomPick = (arr: (bigint | number)[], times: number) => {
  const randomArr = [...Array(arr.length).keys()].map((_, idx) => arr[idx]);
  for (let i = 0; i < randomArr.length; i++) {
    const randomPos = Math.floor(Math.random() * arr.length - 1);

    const temp = randomArr[i];
    randomArr[i] = randomArr[randomPos];
    randomArr[randomPos] = temp;
  }

  return randomArr.filter((_, idx) => idx < times);
};

async function main() {
  const usersName = [...Array(USERS_NUMBER).keys()].map(() => ({
    id: generateId(),
    name: faker.person.firstName(),
  }));
  const chatRoomsName = [...Array(ROOMS_NUMBER).keys()].map(() => ({
    id: generateId(),
    name: faker.music.songName(),
  }));

  await prisma.users.createMany({ data: [...usersName] });
  console.log('gen users data successfully!');
  const users = await prisma.users.findMany();

  await prisma.chatRoom.createMany({
    data: [...chatRoomsName],
  });
  console.log('gen rooms data successfully!');
  const rooms = await prisma.chatRoom.findMany();

  {
    const roomIds = randomPick(
      rooms.map((r) => r.id),
      Math.floor(ROOMS_NUMBER * 0.75)
    );
    const userIds = randomPick(
      users.map((u) => u.id),
      Math.floor(USERS_NUMBER * 0.75)
    );
    const createManyDataInput: {
      id: bigint;
      room_id: bigint;
      user_id: bigint;
    }[] = [];
    for (const uId of userIds) {
      for (const rId of roomIds) {
        createManyDataInput.push({
          id: generateId(),
          user_id: BigInt(uId),
          room_id: BigInt(rId),
        });
      }
    }
    await prisma.usersAndChatRooms.createMany({
      data: [...createManyDataInput],
    });
    console.log('Assign user into chat rooms successfully!');
  }

  // Each user go inside joined room and talk shit
  {
    const usersAndRooms = await prisma.usersAndChatRooms.findMany();

    const messagesDataInput = [];

    for (const i of usersAndRooms) {
      const contents: string[] = [];
      for (let i = 0; i < Math.floor(Math.random() * 2); i++) {
        contents.push(faker.word.words(Math.floor(Math.random() * 20)));
      }

      for (const content of contents) {
        messagesDataInput.push({
          id: generateId(),
          user_id: i.user_id,
          room_id: i.room_id,
          type: 'text',
          content,
        });
      }
    }

    console.log('faking messages data');
    await prisma.message.createMany({ data: [...messagesDataInput] });
    console.log('Successfully!');
  }
}

main()
  .then(async () => {
    await prisma.$disconnect();
  })
  .catch(async (e) => {
    console.error(e);
    await prisma.$disconnect();
    process.exit(1);
  });
