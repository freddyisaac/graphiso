#include <rav1e.h>
#include <stdio.h>
#include <inttypes.h>

#include "trampoline.h"

RAV1E *new_rav1e()
{
	return (RAV1E*) calloc(sizeof(RAV1E), 1);
}

void t_rav1e_config_default(RAV1E *r)
{
	r->rac = rav1e_config_default();
}

void t_rav1e_clean(RAV1E *r)
{
	if (!r) return;
	if (r->raf) rav1e_frame_unref(r->raf);
	if (r->rax) rav1e_context_unref(r->rax);
	if (r->rac) rav1e_config_unref(r->rac);
	free(r);
}

int t_rav1e_simple_setup(RAV1E *r)
{
	int ret;
	if (!r) return -1;
    if (!r->rac) {
        printf("Unable to initialize\n");
    }
	ret = rav1e_config_parse_int(r->rac, "width", 64);
    if (ret < 0) {
        printf("Unable to configure width\n");
		return ret;
    }

    ret = rav1e_config_parse_int(r->rac, "height", 96);
    if (ret < 0) {
        printf("Unable to configure height\n");
		return ret;
    }

    ret = rav1e_config_parse_int(r->rac, "speed", 9);
    if (ret < 0) {
        printf("Unable to configure speed\n");
		return ret;
    }

    ret = rav1e_config_set_color_description(r->rac, 2, 2, 2);
    if (ret < 0) {
        printf("Unable to configure color properties\n");
		return ret;
    }
	return 0;
}

int t_rav1e_context_and_frame(RAV1E *r)
{
	if (!r) return -1;
	r->rax = rav1e_context_new(r->rac);
	if (!r->rax) {
 		printf("Unable to allocate a new context\n");
		return -1;
	}

	r->raf = rav1e_frame_new(r->rax);
	if (!r->raf) {
		printf("Unable to allocate a new frame\n");
		return -1;
	}
	return 0;
}

int t_rav1e_simple_chromaticity(RAV1E *r)
{
	if (!r) return -1;

	int ret;
    RaChromaticityPoint primaries[] = {
        { .x = 0.68 * (1 << 16),  .y = 0.32 * (1 << 16) },
        { .x = 0.265 * (1 << 16), .y = 0.69 * (1 << 16) },
        { .x = 0.15 * (1 << 16),  .y = 0.06 * (1 << 16) },
    };
    RaChromaticityPoint wp = { .x = 0.31268 * (1 << 16), .y = 0.329 * (1 << 16) };
    ret = rav1e_config_set_mastering_display(r->rac, primaries, wp, 1000 * (1 << 8), 0 * (1 << 14));
    if (ret < 0) {
        printf("Unable to configure mastering display\n");
		return ret;
    }

    ret = rav1e_config_set_content_light(r->rac, 1000, 0);
    if (ret < 0) {
        printf("Unable to configure mastering display\n");
		return ret;
    }
	
	return 0;
}

int t_rav1e_send(RAV1E *r)
{
	if (!r) return -1;
	int ret;
	int limit = 30;
	printf("Encoding %d frames\n", limit);

	int i;
	for (i=0;i<limit;i++)
	{
		printf("sending frame\n");
		ret = rav1e_send_frame(r->rax, r->raf);
		if (ret < 0)
		{
			printf("unable to send frame %d\n");
			return -1;
		}else{
			if (ret > 0) {
				printf("unable to append frame %d to internal queue\n", i);
			}
		}
	}

	rav1e_send_frame(r->rax, NULL);

	ret = 0;
	for (i=0;i<limit+5;)
	{
		RaPacket *p;
		ret = rav1e_receive_packet(r->rax, &p);
		if (ret < 0)
		{
			return ret;
		}else{
			if (ret==0)
			{
				printf("packet %"PRIu64"\n", p->input_frameno);
				rav1e_packet_unref(p);
			}else{
				if (ret == RA_ENCODER_STATUS_LIMIT_REACHED)
				{
					printf("Limit reached\n");
					break;
				}
			}
		}
	}
}

void F()
{
	RaConfig *rac = rav1e_config_default();
	RaFrame *f = NULL;
    RaContext *rax = NULL;
    int ret = -1;

    if (!rac) {
        printf("Unable to initialize\n");
        goto clean;
    }

   ret = rav1e_config_parse_int(rac, "width", 64);
    if (ret < 0) {
        printf("Unable to configure width\n");
        goto clean;
    }

    ret = rav1e_config_parse_int(rac, "height", 96);
    if (ret < 0) {
        printf("Unable to configure height\n");
        goto clean;
    }

    ret = rav1e_config_parse_int(rac, "speed", 9);
    if (ret < 0) {
        printf("Unable to configure speed\n");
        goto clean;
    }

    ret = rav1e_config_set_color_description(rac, 2, 2, 2);
    if (ret < 0) {
        printf("Unable to configure color properties\n");
        goto clean;
    }

   rax = rav1e_context_new(rac);
   if (!rax) {
       printf("Unable to allocate a new context\n");
       goto clean;
   }

   f = rav1e_frame_new(rax);
   if (!f) {
       printf("Unable to allocate a new frame\n");
       goto clean;
   }


clean:
	if (!f) rav1e_frame_unref(f);
	if (!rax) rav1e_context_unref(rax);
	if (!rac) rav1e_config_unref(rac);
}

