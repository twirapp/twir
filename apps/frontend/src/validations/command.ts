import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import yup from 'yup';

const createCommandValidation = (command: UpdateOrCreateCommandDto) => {
  return yup.object({
    name: yup.string().required()
      .test(
        'unique-name',
        (d) => `${d.path} already used for other command`,
        (v) => {
          if (!v) return false;
          const otherCommands = commands.value?.filter(c => c.id !== command.id);

          if (otherCommands?.some(c => c.name === v)) {
            return false;
          }

          if (otherCommands?.some(c => c.aliases?.some(aliases => aliases.includes(v)))) {
            return false;
          }

          return true;
        },
      ),
    cooldown: yup.number().optional().min(5),
    cooldownType: yup.string().optional().test((v) => v ? Object.values<string>(perms).includes(v) : false),
    responses: yup.array<CommandType['responses']>()
      .of(yup.object().shape({
        text: yup.string().required(),
      })),
  });
};