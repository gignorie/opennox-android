#version 150 core

in vec2 Texcoord;

out vec4 color;

uniform sampler2D tex;
uniform float gamma = 1.0;

void main()
{
    color = texture(tex, Texcoord);
    color.rgb = pow(color.rgb, vec3(1.0/gamma));
}
